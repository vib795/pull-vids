import { NextRequest, NextResponse } from 'next/server'
import ytdl from 'ytdl-core'

export async function POST(request: NextRequest) {
  try {
    const { url, quality, audioOnly } = await request.json()

    // Validate URL
    if (!url) {
      return NextResponse.json(
        { error: 'URL is required' },
        { status: 400 }
      )
    }

    // For now, we only support YouTube due to serverless limitations
    // Other platforms would require yt-dlp which needs a binary
    if (!ytdl.validateURL(url)) {
      return NextResponse.json(
        { error: 'Invalid YouTube URL. For other platforms, please use the CLI version.' },
        { status: 400 }
      )
    }

    // Get video info
    const info = await ytdl.getInfo(url)
    const title = info.videoDetails.title

    // Choose format based on quality and audioOnly settings
    let format
    if (audioOnly) {
      format = ytdl.chooseFormat(info.formats, { quality: 'highestaudio' })
    } else {
      // Map quality to itag or use highest
      const qualityMap: { [key: string]: string } = {
        '2160p': '137',  // 4K video
        '1440p': '264',  // 2K video
        '1080p': '137',  // 1080p video
        '720p': '136',   // 720p video
        '480p': '135',   // 480p video
        '360p': '134',   // 360p video
        'best': 'highest'
      }

      const selectedQuality = qualityMap[quality] || 'highest'
      format = ytdl.chooseFormat(info.formats, {
        quality: selectedQuality,
        filter: audioOnly ? 'audioonly' : 'videoandaudio'
      })
    }

    if (!format) {
      return NextResponse.json(
        { error: 'Could not find suitable format' },
        { status: 500 }
      )
    }

    // Return download URL
    // Note: In production, you might want to stream this or use a different approach
    // due to Vercel's serverless function limitations
    return NextResponse.json({
      downloadUrl: format.url,
      title: title,
      format: format.container,
      quality: format.qualityLabel || 'audio'
    })
  } catch (error) {
    console.error('Download error:', error)

    if (error instanceof Error) {
      return NextResponse.json(
        { error: error.message },
        { status: 500 }
      )
    }

    return NextResponse.json(
      { error: 'Failed to download video' },
      { status: 500 }
    )
  }
}

// Info endpoint to get video details
export async function GET(request: NextRequest) {
  try {
    const url = request.nextUrl.searchParams.get('url')

    if (!url) {
      return NextResponse.json(
        { error: 'URL parameter is required' },
        { status: 400 }
      )
    }

    if (!ytdl.validateURL(url)) {
      return NextResponse.json(
        { error: 'Invalid YouTube URL' },
        { status: 400 }
      )
    }

    const info = await ytdl.getInfo(url)

    return NextResponse.json({
      title: info.videoDetails.title,
      duration: info.videoDetails.lengthSeconds,
      thumbnail: info.videoDetails.thumbnails[0]?.url,
      author: info.videoDetails.author.name,
      views: info.videoDetails.viewCount,
      formats: info.formats.map(f => ({
        quality: f.qualityLabel,
        container: f.container,
        hasAudio: f.hasAudio,
        hasVideo: f.hasVideo,
      }))
    })
  } catch (error) {
    console.error('Info error:', error)
    return NextResponse.json(
      { error: 'Failed to fetch video info' },
      { status: 500 }
    )
  }
}
