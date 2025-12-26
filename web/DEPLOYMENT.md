# Deploying pull-vids Web App

Complete guide to deploy the web version of pull-vids for free.

## 🚀 Quick Deploy (Vercel)

### One-Click Deploy

Click this button to deploy instantly:

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/vib795/pull-vids/tree/main/web)

### Manual Deploy

#### Step 1: Prerequisites

- GitHub account
- Vercel account (sign up at https://vercel.com)
- This repository pushed to GitHub

#### Step 2: Import Project

1. Go to https://vercel.com/new
2. Click "Import Git Repository"
3. Select your `pull-vids` repository
4. **Root Directory:** Set to `web`
5. **Framework Preset:** Next.js (auto-detected)
6. Click "Deploy"

#### Step 3: Configure (Optional)

If you want custom domain:
1. Go to project settings in Vercel
2. Navigate to "Domains"
3. Add your custom domain
4. Follow DNS configuration steps

That's it! Your app is live! 🎉

---

## 🌍 Alternative Free Hosting Options

### Netlify

1. Go to https://app.netlify.com
2. Click "Add new site" → "Import an existing project"
3. Connect to GitHub
4. Select repository
5. **Base directory:** `web`
6. **Build command:** `npm run build`
7. **Publish directory:** `.next`
8. Deploy!

### Cloudflare Pages

1. Go to https://dash.cloudflare.com
2. Navigate to "Pages"
3. Click "Create a project"
4. Connect GitHub repository
5. **Build command:** `npm run build`
6. **Build output directory:** `.next`
7. Deploy!

### Railway

1. Go to https://railway.app
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Select repository
5. **Root directory:** `web`
6. Railway auto-detects Next.js
7. Deploy!

---

## 🔧 Configuration

### Environment Variables

If you need to add environment variables:

**Vercel:**
1. Go to Project Settings
2. Navigate to "Environment Variables"
3. Add variables

**Common variables:**
```env
# Add any future API keys here
NEXT_PUBLIC_API_URL=https://your-domain.com
```

### Custom Domain

**Vercel:**
1. Project Settings → Domains
2. Add domain
3. Configure DNS:
   - Type: A
   - Name: @
   - Value: 76.76.21.21
   - Type: CNAME
   - Name: www
   - Value: cname.vercel-dns.com

**Netlify:**
1. Site Settings → Domain management
2. Add custom domain
3. Follow DNS instructions

---

## 📊 Monitoring

### Vercel Analytics

Enable free analytics:
1. Go to project in Vercel
2. Navigate to "Analytics"
3. Enable Analytics
4. View real-time stats

### Error Tracking

View errors and logs:
1. Vercel Dashboard → Your Project
2. Navigate to "Logs"
3. Filter by errors

---

## 🚦 Testing Before Deploy

### Local Testing

```bash
cd web

# Install dependencies
npm install

# Run dev server
npm run dev

# Test at http://localhost:3000
```

### Production Build Locally

```bash
# Build production version
npm run build

# Test production build
npm start

# Should work perfectly at http://localhost:3000
```

### Check Lighthouse Score

1. Open Chrome DevTools
2. Go to "Lighthouse" tab
3. Run audit
4. Aim for 90+ scores

---

## 🔒 Security

### Rate Limiting

Consider adding rate limiting for production:

```typescript
// lib/rate-limit.ts
import { NextRequest } from 'next/server'

export function rateLimit(req: NextRequest) {
  // Implement rate limiting logic
  // Use Vercel Edge Config or Upstash Redis
}
```

### CORS

Configure CORS if needed:

```typescript
// next.config.js
module.exports = {
  async headers() {
    return [
      {
        source: '/api/:path*',
        headers: [
          { key: 'Access-Control-Allow-Origin', value: '*' },
        ],
      },
    ]
  },
}
```

---

## 📈 Performance Tips

### 1. Enable Caching

```typescript
// app/api/download/route.ts
export const revalidate = 3600 // Cache for 1 hour
```

### 2. Optimize Images

Use Next.js Image component:
```tsx
import Image from 'next/image'

<Image src="/icon.png" width={32} height={32} alt="Icon" />
```

### 3. Enable Compression

Vercel automatically enables compression.

For other hosts, add:
```javascript
// next.config.js
module.exports = {
  compress: true,
}
```

---

## 🐛 Troubleshooting

### Build Fails

**Issue:** npm install fails
```bash
# Clear cache
rm -rf node_modules package-lock.json
npm install
```

**Issue:** TypeScript errors
```bash
# Run type check
npm run type-check

# Fix errors and rebuild
```

### Runtime Errors

**Issue:** API route not working
- Check Vercel function logs
- Ensure environment variables are set
- Verify API route path matches

**Issue:** Downloads fail
- Check ytdl-core version compatibility
- Verify YouTube URL is valid
- Check Vercel function timeout (max 30s on free tier)

### Performance Issues

**Issue:** Slow loading
- Enable Vercel Edge Network
- Use Next.js Image optimization
- Implement proper caching

---

## 📋 Pre-Deployment Checklist

- [ ] Test locally with `npm run dev`
- [ ] Build succeeds with `npm run build`
- [ ] All pages load correctly
- [ ] API routes work
- [ ] Forms submit properly
- [ ] Error handling works
- [ ] Mobile responsive
- [ ] Lighthouse score > 90
- [ ] Security headers configured
- [ ] Analytics enabled (optional)

---

## 🌟 Post-Deployment

### 1. Test Live Site

- Visit your deployed URL
- Test video download
- Try different qualities
- Test audio extraction
- Check mobile view

### 2. Set Up Monitoring

- Enable Vercel Analytics
- Set up error alerts
- Monitor function usage

### 3. Update Repository

Update main README with web app link:

```markdown
## 🌐 Web Version

Try the web version: https://your-app.vercel.app

No installation needed - just paste a URL and download!
```

### 4. Share

- Post on Twitter/X
- Share on LinkedIn
- Submit to Product Hunt
- Post on Reddit (r/selfhosted, r/InternetIsBeautiful)

---

## 💰 Costs

### Vercel Free Tier

- **Bandwidth:** 100 GB/month
- **Function Executions:** 100 GB-hours
- **Build Time:** 100 hours/month
- **Perfect for:** Personal projects, demos

### Upgrade If Needed

If you exceed free tier:
- **Hobby ($20/month):** Unlimited bandwidth
- **Pro ($40/month):** Team features, more functions

---

## 🎓 Learning Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [Vercel Documentation](https://vercel.com/docs)
- [Next.js Deployment](https://nextjs.org/docs/deployment)
- [Vercel Edge Functions](https://vercel.com/docs/functions/edge-functions)

---

**Need help?** Open an issue at https://github.com/vib795/pull-vids/issues

Good luck with your deployment! 🚀
