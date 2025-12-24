#!/usr/bin/env python3
"""
pull-vids - A free CLI tool for downloading YouTube videos and audio.
Inspired by Downie and PullTube, but completely free and open-source.
"""

import argparse
import sys
import os
from pathlib import Path
import yt_dlp
from colorama import init, Fore, Style

# Initialize colorama for cross-platform colored output
init(autoreset=True)

__version__ = "0.1.0"


def print_banner():
    """Print the application banner."""
    banner = f"""
{Fore.CYAN}╔═══════════════════════════════════════╗
║          pull-vids v{__version__}           ║
║   Free YouTube Video Downloader CLI   ║
╚═══════════════════════════════════════╝{Style.RESET_ALL}
"""
    print(banner)


def get_format_string(quality, audio_only=False):
    """
    Get the format string for yt-dlp based on quality selection.

    Args:
        quality: Quality setting (best, high, medium, low, or specific like 1080p, 720p)
        audio_only: Whether to download audio only

    Returns:
        Format string for yt-dlp
    """
    if audio_only:
        return 'bestaudio/best'

    quality_map = {
        'best': 'bestvideo+bestaudio/best',
        'high': 'bestvideo[height<=1080]+bestaudio/best[height<=1080]',
        'medium': 'bestvideo[height<=720]+bestaudio/best[height<=720]',
        'low': 'bestvideo[height<=480]+bestaudio/best[height<=480]',
        '2160p': 'bestvideo[height<=2160]+bestaudio/best[height<=2160]',
        '1440p': 'bestvideo[height<=1440]+bestaudio/best[height<=1440]',
        '1080p': 'bestvideo[height<=1080]+bestaudio/best[height<=1080]',
        '720p': 'bestvideo[height<=720]+bestaudio/best[height<=720]',
        '480p': 'bestvideo[height<=480]+bestaudio/best[height<=480]',
        '360p': 'bestvideo[height<=360]+bestaudio/best[height<=360]',
    }

    return quality_map.get(quality, 'bestvideo+bestaudio/best')


def progress_hook(d):
    """Hook to display download progress."""
    if d['status'] == 'downloading':
        try:
            percent = d.get('_percent_str', 'N/A')
            speed = d.get('_speed_str', 'N/A')
            eta = d.get('_eta_str', 'N/A')

            # Clear line and print progress
            print(f'\r{Fore.GREEN}Downloading: {percent} | Speed: {speed} | ETA: {eta}{Style.RESET_ALL}', end='', flush=True)
        except:
            pass
    elif d['status'] == 'finished':
        print(f'\n{Fore.GREEN}✓ Download complete! Processing...{Style.RESET_ALL}')


def download_video(url, output_dir, quality='best', audio_only=False, playlist=False, format_type=None):
    """
    Download a YouTube video or playlist.

    Args:
        url: YouTube URL to download
        output_dir: Directory to save downloads
        quality: Video quality to download
        audio_only: Download audio only
        playlist: Allow playlist downloads
        format_type: Output format (mp4, mkv, mp3, m4a, etc.)
    """
    # Ensure output directory exists
    output_path = Path(output_dir).expanduser().resolve()
    output_path.mkdir(parents=True, exist_ok=True)

    # Build output template
    if playlist:
        output_template = str(output_path / '%(playlist)s/%(playlist_index)s - %(title)s.%(ext)s')
    else:
        output_template = str(output_path / '%(title)s.%(ext)s')

    # Configure yt-dlp options
    ydl_opts = {
        'format': get_format_string(quality, audio_only),
        'outtmpl': output_template,
        'progress_hooks': [progress_hook],
        'quiet': False,
        'no_warnings': False,
        'noplaylist': not playlist,
    }

    # Add audio-specific options
    if audio_only:
        preferred_format = format_type if format_type else 'mp3'
        ydl_opts['postprocessors'] = [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': preferred_format,
            'preferredquality': '192',
        }]
    elif format_type:
        ydl_opts['merge_output_format'] = format_type

    try:
        print(f"{Fore.CYAN}Starting download from: {url}{Style.RESET_ALL}")
        print(f"{Fore.CYAN}Output directory: {output_path}{Style.RESET_ALL}")

        if audio_only:
            print(f"{Fore.YELLOW}Mode: Audio only{Style.RESET_ALL}")
        else:
            print(f"{Fore.YELLOW}Quality: {quality}{Style.RESET_ALL}")

        print()

        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            # Extract info first to show what we're downloading
            info = ydl.extract_info(url, download=False)

            if 'entries' in info and playlist:
                print(f"{Fore.MAGENTA}Playlist detected: {info.get('title', 'Unknown')}{Style.RESET_ALL}")
                print(f"{Fore.MAGENTA}Videos in playlist: {len(info['entries'])}{Style.RESET_ALL}\n")
            else:
                title = info.get('title', 'Unknown')
                duration = info.get('duration', 0)
                mins, secs = divmod(duration, 60)
                print(f"{Fore.MAGENTA}Title: {title}{Style.RESET_ALL}")
                print(f"{Fore.MAGENTA}Duration: {mins}m {secs}s{Style.RESET_ALL}\n")

            # Perform the download
            ydl.download([url])

        print(f"\n{Fore.GREEN}{'=' * 50}{Style.RESET_ALL}")
        print(f"{Fore.GREEN}✓ All downloads completed successfully!{Style.RESET_ALL}")
        print(f"{Fore.GREEN}{'=' * 50}{Style.RESET_ALL}\n")

        return True

    except yt_dlp.utils.DownloadError as e:
        print(f"\n{Fore.RED}✗ Download error: {str(e)}{Style.RESET_ALL}", file=sys.stderr)
        return False
    except KeyboardInterrupt:
        print(f"\n{Fore.YELLOW}✗ Download cancelled by user{Style.RESET_ALL}")
        return False
    except Exception as e:
        print(f"\n{Fore.RED}✗ Unexpected error: {str(e)}{Style.RESET_ALL}", file=sys.stderr)
        return False


def main():
    """Main entry point for the CLI."""
    parser = argparse.ArgumentParser(
        description='pull-vids - Free YouTube video and audio downloader',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Download video in best quality
  pull-vids "https://www.youtube.com/watch?v=VIDEO_ID"

  # Download audio only as MP3
  pull-vids -a "https://www.youtube.com/watch?v=VIDEO_ID"

  # Download in 720p quality to specific directory
  pull-vids -q 720p -o ~/Videos "https://www.youtube.com/watch?v=VIDEO_ID"

  # Download entire playlist
  pull-vids -p "https://www.youtube.com/playlist?list=PLAYLIST_ID"

  # Download audio from playlist
  pull-vids -a -p "https://www.youtube.com/playlist?list=PLAYLIST_ID"
        """
    )

    parser.add_argument(
        'url',
        help='YouTube video or playlist URL'
    )

    parser.add_argument(
        '-o', '--output',
        default='~/Downloads/pull-vids',
        help='Output directory (default: ~/Downloads/pull-vids)'
    )

    parser.add_argument(
        '-q', '--quality',
        choices=['best', 'high', 'medium', 'low', '2160p', '1440p', '1080p', '720p', '480p', '360p'],
        default='best',
        help='Video quality (default: best)'
    )

    parser.add_argument(
        '-a', '--audio-only',
        action='store_true',
        help='Download audio only'
    )

    parser.add_argument(
        '-p', '--playlist',
        action='store_true',
        help='Download entire playlist'
    )

    parser.add_argument(
        '-f', '--format',
        help='Output format (mp4, mkv, mp3, m4a, etc.)'
    )

    parser.add_argument(
        '-v', '--version',
        action='version',
        version=f'pull-vids {__version__}'
    )

    parser.add_argument(
        '--no-banner',
        action='store_true',
        help='Don\'t show the banner'
    )

    args = parser.parse_args()

    # Show banner unless disabled
    if not args.no_banner:
        print_banner()

    # Validate URL
    if not args.url or not ('youtube.com' in args.url or 'youtu.be' in args.url):
        print(f"{Fore.RED}✗ Invalid YouTube URL{Style.RESET_ALL}", file=sys.stderr)
        return 1

    # Download the video
    success = download_video(
        url=args.url,
        output_dir=args.output,
        quality=args.quality,
        audio_only=args.audio_only,
        playlist=args.playlist,
        format_type=args.format
    )

    return 0 if success else 1


if __name__ == '__main__':
    sys.exit(main())
