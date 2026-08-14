#!/usr/bin/env python3
"""Render assets/icon.png, the package icon referenced by the Chocolatey nuspec.

Standard library only: no PIL, ImageMagick or SVG rasteriser is installed on
the machines that build this, and pulling one in just to redraw a static mark
is not worth the dependency. Edges come from supersampling rather than
analytic anti-aliasing, which is slower but keeps the geometry obvious.

    python3 assets/make-icon.py

The mark is a download arrow whose head is also a play triangle, so one shape
reads as both "download" and "video". Rendered at 256px because Chocolatey
displays at 128 and a clean 2x downscale beats a blurry upscale.
"""

import os
import struct
import zlib

SIZE = 256
SS = 4  # supersamples per axis

# Cyan echoes the banner the CLI prints, so the icon and the tool agree.
GRAD_TOP = (56, 224, 245)
GRAD_BOTTOM = (10, 116, 172)
FG = (255, 255, 255)

CORNER_RADIUS = 56.0

# Arrow stem.
STEM_X0, STEM_X1 = 111.0, 145.0
STEM_Y0, STEM_Y1 = 48.0, 118.0

# Arrowhead: a down-pointing triangle that doubles as a play glyph.
HEAD_X0, HEAD_X1 = 68.0, 188.0
HEAD_Y_TOP, HEAD_Y_APEX = 106.0, 178.0

# Base bar the arrow points into.
BAR_X0, BAR_X1 = 68.0, 188.0
BAR_Y0, BAR_Y1 = 196.0, 214.0
BAR_RADIUS = 9.0


def in_rounded_rect(x, y, x0, y0, x1, y1, r):
    """Point-in-rounded-rectangle, by clamping to the inner corner circle."""
    if not (x0 <= x <= x1 and y0 <= y <= y1):
        return False
    cx = min(max(x, x0 + r), x1 - r)
    cy = min(max(y, y0 + r), y1 - r)
    dx, dy = x - cx, y - cy
    if dx == 0.0 or dy == 0.0:
        return True
    return dx * dx + dy * dy <= r * r


def in_down_triangle(x, y, x0, x1, y_top, y_apex):
    """Point-in-triangle for an isoceles triangle pointing down."""
    if y < y_top or y > y_apex:
        return False
    # Half-width shrinks linearly from the top edge to the apex.
    t = (y - y_top) / (y_apex - y_top)
    half = (x1 - x0) / 2.0 * (1.0 - t)
    mid = (x0 + x1) / 2.0
    return mid - half <= x <= mid + half


def in_glyph(x, y):
    return (
        in_rounded_rect(x, y, STEM_X0, STEM_Y0, STEM_X1, STEM_Y1, 6.0)
        or in_down_triangle(x, y, HEAD_X0, HEAD_X1, HEAD_Y_TOP, HEAD_Y_APEX)
        or in_rounded_rect(x, y, BAR_X0, BAR_Y0, BAR_X1, BAR_Y1, BAR_RADIUS)
    )


def render():
    rows = []
    step = 1.0 / SS
    offset = step / 2.0
    total = float(SS * SS)

    for py in range(SIZE):
        # Vertical gradient, sampled once per row.
        t = py / float(SIZE - 1)
        base = tuple(
            int(round(GRAD_TOP[i] + (GRAD_BOTTOM[i] - GRAD_TOP[i]) * t))
            for i in range(3)
        )

        row = bytearray()
        for px in range(SIZE):
            bg_hits = 0
            fg_hits = 0
            for sy in range(SS):
                y = py + offset + sy * step
                for sx in range(SS):
                    x = px + offset + sx * step
                    if not in_rounded_rect(x, y, 0.0, 0.0, SIZE, SIZE,
                                           CORNER_RADIUS):
                        continue
                    bg_hits += 1
                    if in_glyph(x, y):
                        fg_hits += 1

            if bg_hits == 0:
                row += b"\x00\x00\x00\x00"
                continue

            # Glyph coverage is measured against the samples that landed inside
            # the tile, not against all samples. Dividing by the total instead
            # would dim the glyph wherever it meets the rounded edge.
            cover = fg_hits / float(bg_hits)
            colour = tuple(
                int(round(base[i] + (FG[i] - base[i]) * cover)) for i in range(3)
            )
            alpha = int(round(255.0 * bg_hits / total))
            row += bytes((colour[0], colour[1], colour[2], alpha))
        rows.append(bytes(row))
    return rows


def write_png(path, rows):
    raw = b"".join(b"\x00" + r for r in rows)

    def chunk(tag, data):
        body = tag + data
        return (struct.pack(">I", len(data)) + body
                + struct.pack(">I", zlib.crc32(body) & 0xFFFFFFFF))

    ihdr = struct.pack(">IIBBBBB", SIZE, SIZE, 8, 6, 0, 0, 0)
    png = (b"\x89PNG\r\n\x1a\n"
           + chunk(b"IHDR", ihdr)
           + chunk(b"IDAT", zlib.compress(raw, 9))
           + chunk(b"IEND", b""))
    with open(path, "wb") as fh:
        fh.write(png)


if __name__ == "__main__":
    out = os.path.join(os.path.dirname(os.path.abspath(__file__)), "icon.png")
    write_png(out, render())
    print("wrote %s (%d bytes)" % (out, os.path.getsize(out)))
