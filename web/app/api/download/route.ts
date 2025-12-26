import { NextRequest, NextResponse } from 'next/server'

export async function POST(request: NextRequest) {
  try {
    const { url } = await request.json()

    // Validate URL
    if (!url) {
      return NextResponse.json(
        { error: 'URL is required' },
        { status: 400 }
      )
    }

    // Important message about limitations
    return NextResponse.json(
      {
        error: 'Direct downloads are currently unavailable due to platform restrictions.',
        message: 'For reliable downloads, please use the CLI version:',
        cli_install: 'brew tap vib795/tap && brew install pull-vids',
        cli_usage: `pull-vids '${url}'`,
        alternative: 'Or visit: https://github.com/vib795/pull-vids',
        reason: 'YouTube and other platforms block requests from serverless functions. The CLI version uses yt-dlp which is more reliable and supports 1000+ platforms.',
      },
      { status: 503 }
    )
  } catch (error) {
    console.error('Download error:', error)

    return NextResponse.json(
      {
        error: 'Service temporarily unavailable',
        message: 'Please use the CLI version for reliable downloads',
        install: 'brew tap vib795/tap && brew install pull-vids'
      },
      { status: 503 }
    )
  }
}

// Info endpoint - also disabled for now
export async function GET(request: NextRequest) {
  return NextResponse.json(
    {
      message: 'Video downloads via web interface are currently unavailable',
      reason: 'Platform restrictions on serverless functions',
      solution: 'Use the CLI version for full functionality',
      install: {
        macos: 'brew tap vib795/tap && brew install pull-vids',
        linux: 'See: https://github.com/vib795/pull-vids#installation',
        windows: 'See: https://github.com/vib795/pull-vids#installation',
      },
      features: 'CLI supports 1000+ platforms, all qualities, playlists, and more'
    }
  )
}
