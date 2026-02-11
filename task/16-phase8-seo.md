# Phase 9: SEO & Discoverability

## Goal
Optimize Texly.Chat for search engines and AI crawlers to maximize organic traffic and discoverability.

---

## Backend Tasks

### Step 1: SEO-Friendly Routes & Metadata

#### 1.1 Server-Side Rendering (SSR) Setup
- [ ] Update TanStack Start configuration for SSR:
  - [ ] Enable SSR for marketing pages (landing, pricing, about)
  - [ ] Keep dashboard as client-side rendered (CSR)
  - [ ] Ensure meta tags rendered on server for crawlers

#### 1.2 Dynamic Meta Tags
- [ ] Create `ui/src/lib/seo.ts`:
  ```ts
  export const defaultSEO = {
    title: "Texly - Build AI Chatbots Trained on Your Data",
    description: "Create custom AI chatbots with RAG-powered responses in minutes. Upload PDFs, Excel, URLs, and text to train your bot. Embed anywhere with one line of code.",
    keywords: "AI chatbot, RAG, chatbot builder, custom chatbot, embed chatbot, AI customer support",
    ogImage: "/og-image.png",
    twitterCard: "summary_large_image",
  };
  
  export const pageSEO = {
    pricing: {
      title: "Pricing - Texly AI Chatbots",
      description: "Start free or upgrade to Pro for $20/month with pay-as-you-go pricing. Enterprise plans available for high-volume users.",
    },
    about: {
      title: "About - Texly AI Chatbots",
      description: "Learn about Texly's mission to make AI chatbots accessible to everyone with RAG-powered accuracy.",
    },
    // Add more pages...
  };
  ```

#### 1.3 Meta Tag Component
- [ ] Create `ui/src/components/seo/meta-tags.tsx`:
  ```tsx
  import { Head } from '@tanstack/react-start';
  
  export function MetaTags({ title, description, keywords, ogImage, canonical }) {
    return (
      <Head>
        <title>{title}</title>
        <meta name="description" content={description} />
        <meta name="keywords" content={keywords} />
        <link rel="canonical" href={canonical} />
        
        {/* Open Graph */}
        <meta property="og:title" content={title} />
        <meta property="og:description" content={description} />
        <meta property="og:image" content={ogImage} />
        <meta property="og:type" content="website" />
        
        {/* Twitter Card */}
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content={title} />
        <meta name="twitter:description" content={description} />
        <meta name="twitter:image" content={ogImage} />
      </Head>
    );
  }
  ```

#### 1.4 Apply Meta Tags to Pages
- [ ] Update `ui/src/routes/index.tsx` (landing):
  - [ ] Add `<MetaTags {...defaultSEO} />`
- [ ] Update `ui/src/routes/pricing/index.tsx`:
  - [ ] Add `<MetaTags {...pageSEO.pricing} />`
- [ ] Repeat for all public pages

---

### Step 2: Sitemap Generation

#### 2.1 Create Sitemap Endpoint
- [ ] Create `internal/handlers/seo/sitemap_handler.go`:
  ```go
  func (h *SEOHandler) GenerateSitemap(c *gin.Context) {
      sitemap := `<?xml version="1.0" encoding="UTF-8"?>
      <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
          <url>
              <loc>https://texly.chat/</loc>
              <lastmod>2025-02-11</lastmod>
              <changefreq>weekly</changefreq>
              <priority>1.0</priority>
          </url>
          <url>
              <loc>https://texly.chat/pricing</loc>
              <lastmod>2025-02-11</lastmod>
              <changefreq>monthly</changefreq>
              <priority>0.8</priority>
          </url>
          <url>
              <loc>https://texly.chat/about</loc>
              <lastmod>2025-02-11</lastmod>
              <changefreq>monthly</changefreq>
              <priority>0.7</priority>
          </url>
          <!-- Add more static pages -->
      </urlset>`
      
      c.Header("Content-Type", "application/xml")
      c.String(200, sitemap)
  }
  ```

#### 2.2 Dynamic Sitemap (Blog Posts)
- [ ] If blog is implemented, dynamically generate URLs:
  - [ ] Query all published blog posts
  - [ ] Add to sitemap XML
  - [ ] Update `lastmod` based on post update date

#### 2.3 Register Sitemap Route
- [ ] Add to `internal/server/server.go`:
  ```go
  router.GET("/sitemap.xml", seoHandler.GenerateSitemap)
  ```

---

### Step 3: Robots.txt

#### 3.1 Create Robots.txt Endpoint
- [ ] Create `internal/handlers/seo/robots_handler.go`:
  ```go
  func (h *SEOHandler) GenerateRobotsTxt(c *gin.Context) {
      robotsTxt := `User-agent: *
  Allow: /
  Disallow: /dashboard/
  Disallow: /api/
  
  Sitemap: https://texly.chat/sitemap.xml
  `
      c.Header("Content-Type", "text/plain")
      c.String(200, robotsTxt)
  }
  ```

#### 3.2 Register Robots.txt Route
- [ ] Add to `internal/server/server.go`:
  ```go
  router.GET("/robots.txt", seoHandler.GenerateRobotsTxt)
  ```

---

### Step 4: LLM-Specific Routes

#### 4.1 Create llms.txt
- [ ] Create `internal/handlers/seo/llms_handler.go`:
  ```go
  func (h *SEOHandler) GenerateLLMsTxt(c *gin.Context) {
      llmsTxt := `# Texly.Chat - AI Chatbot Platform
  
  ## About
  Texly is a SaaS platform for building custom AI chatbots trained on your own data using Retrieval-Augmented Generation (RAG).
  
  ## Features
  - RAG-powered chatbots
  - Multi-source support (PDF, Excel, URLs, text)
  - Embeddable widget
  - Real-time chat with streaming
  
  ## Pricing
  - Free tier: 1 bot, 100 messages/month
  - Pro tier: $20/month + pay-as-you-go
  - Enterprise tier: Custom pricing
  
  ## Links
  - Homepage: https://texly.chat/
  - Pricing: https://texly.chat/pricing
  - Documentation: https://texly.chat/docs (if exists)
  - API: https://texly.chat/api/docs (if exists)
  
  ## Contact
  - Email: support@texly.chat
  - Twitter: @texlychat
  `
      c.Header("Content-Type", "text/plain")
      c.String(200, llmsTxt)
  }
  ```

#### 4.2 Create ai.txt
- [ ] Create similar endpoint for `ai.txt`:
  ```go
  func (h *SEOHandler) GenerateAiTxt(c *gin.Context) {
      aiTxt := `# AI Information for Texly.Chat
  
  This is an AI chatbot platform. We use OpenAI's GPT models for chat completions and text-embedding models for RAG.
  
  ## Data Privacy
  - User data is not used to train AI models
  - Data is stored securely and privately
  - Users own their data
  
  ## Allowed AI Access
  User-agent: GPTBot
  Allow: /
  Disallow: /dashboard/
  
  User-agent: ChatGPT-User
  Allow: /
  Disallow: /dashboard/
  `
      c.Header("Content-Type", "text/plain")
      c.String(200, aiTxt)
  }
  ```

#### 4.3 Register Routes
- [ ] Add to server:
  ```go
  router.GET("/llms.txt", seoHandler.GenerateLLMsTxt)
  router.GET("/ai.txt", seoHandler.GenerateAiTxt)
  ```

---

### Step 5: Structured Data (JSON-LD)

#### 5.1 Organization Schema
- [ ] Add to landing page footer or `<Head>`:
  ```html
  <script type="application/ld+json">
  {
    "@context": "https://schema.org",
    "@type": "Organization",
    "name": "Texly",
    "url": "https://texly.chat",
    "logo": "https://texly.chat/logo.png",
    "sameAs": [
      "https://twitter.com/texlychat",
      "https://github.com/yourusername/texly"
    ],
    "contactPoint": {
      "@type": "ContactPoint",
      "email": "support@texly.chat",
      "contactType": "Customer Support"
    }
  }
  </script>
  ```

#### 5.2 Product Schema (Pricing Page)
- [ ] Add to pricing page:
  ```html
  <script type="application/ld+json">
  {
    "@context": "https://schema.org",
    "@type": "Product",
    "name": "Texly Pro",
    "description": "AI chatbot platform with RAG",
    "offers": {
      "@type": "Offer",
      "price": "20.00",
      "priceCurrency": "USD",
      "priceValidUntil": "2026-12-31",
      "availability": "https://schema.org/InStock"
    }
  }
  </script>
  ```

#### 5.3 FAQ Schema (FAQ Page)
- [ ] Add to support/FAQ page:
  ```html
  <script type="application/ld+json">
  {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    "mainEntity": [
      {
        "@type": "Question",
        "name": "How do I create a bot?",
        "acceptedAnswer": {
          "@type": "Answer",
          "text": "Go to Dashboard > Create Bot, add your data sources, and customize."
        }
      }
      // Add more FAQs
    ]
  }
  </script>
  ```

---

## Frontend Tasks

### Step 6: Performance Optimization (Core Web Vitals)

#### 6.1 Image Optimization
- [ ] Use Next.js Image component or similar for lazy loading
- [ ] Optimize images:
  - [ ] Convert to WebP format
  - [ ] Compress with tools like Squoosh or ImageOptim
  - [ ] Set appropriate sizes and srcset
- [ ] Lazy load images below the fold

#### 6.2 Code Splitting
- [ ] Ensure TanStack Router code-splits by route
- [ ] Lazy load heavy components:
  ```tsx
  const HeavyChart = lazy(() => import('./HeavyChart'));
  ```

#### 6.3 Font Optimization
- [ ] Use `font-display: swap` for web fonts
- [ ] Preload critical fonts:
  ```html
  <link rel="preload" href="/fonts/inter.woff2" as="font" type="font/woff2" crossorigin>
  ```
- [ ] Subset fonts to include only used characters

#### 6.4 CSS Optimization
- [ ] TailwindCSS v4 tree-shakes unused styles
- [ ] Inline critical CSS for above-the-fold content
- [ ] Minify CSS in production

#### 6.5 JavaScript Optimization
- [ ] Minimize bundle size:
  - [ ] Use `bun build` or Vite's tree-shaking
  - [ ] Remove unused dependencies
  - [ ] Use dynamic imports for large libraries
- [ ] Defer non-critical scripts

#### 6.6 Lighthouse Audit
- [ ] Run Lighthouse on all pages
- [ ] Target scores:
  - [ ] Performance: >90
  - [ ] Accessibility: >90
  - [ ] Best Practices: >90
  - [ ] SEO: 100
- [ ] Fix all identified issues

---

### Step 7: Open Graph Images

#### 7.1 Create OG Image
- [ ] Design a default OG image (1200x630px):
  - [ ] Texly logo
  - [ ] Tagline: "Build AI Chatbots Trained on Your Data"
  - [ ] Professional design
- [ ] Save as `ui/public/og-image.png`

#### 7.2 Dynamic OG Images (Optional)
- [ ] For blog posts, generate dynamic OG images:
  - [ ] Use library like `@vercel/og` or `node-canvas`
  - [ ] Include post title and author
  - [ ] Endpoint: `/api/og?title=...`

---

### Step 8: Analytics Integration

#### 8.1 Privacy-Friendly Analytics
- [ ] Choose analytics tool:
  - [ ] **Option 1**: Plausible Analytics (privacy-focused, GDPR compliant)
  - [ ] **Option 2**: Google Analytics 4 (with cookie consent)
  - [ ] **Option 3**: Umami (self-hosted, open-source)

#### 8.2 Install Analytics
- [ ] Add analytics script to `ui/src/routes/__root.tsx`:
  - [ ] Only load if user consented to cookies
  - [ ] Track page views
  - [ ] Track custom events (signups, upgrades, bot creations)

#### 8.3 Google Search Console
- [ ] Verify domain ownership:
  - [ ] Add meta tag to `<head>` or upload HTML file
- [ ] Submit sitemap.xml to Search Console
- [ ] Monitor crawl errors and indexing status

#### 8.4 Bing Webmaster Tools
- [ ] Verify domain with Bing
- [ ] Submit sitemap.xml
- [ ] Monitor indexing

---

### Step 9: Social Media Meta Tags

#### 9.1 Twitter Cards
- [ ] Add Twitter-specific meta tags:
  ```html
  <meta name="twitter:card" content="summary_large_image">
  <meta name="twitter:site" content="@texlychat">
  <meta name="twitter:creator" content="@yourhandle">
  <meta name="twitter:title" content="Texly - AI Chatbots">
  <meta name="twitter:description" content="...">
  <meta name="twitter:image" content="https://texly.chat/og-image.png">
  ```

#### 9.2 Open Graph Tags
- [ ] Ensure all pages have OG tags:
  ```html
  <meta property="og:title" content="...">
  <meta property="og:description" content="...">
  <meta property="og:image" content="...">
  <meta property="og:url" content="https://texly.chat">
  <meta property="og:type" content="website">
  <meta property="og:site_name" content="Texly">
  ```

---

### Step 10: Canonical URLs

#### 10.1 Add Canonical Tags
- [ ] Prevent duplicate content:
  ```html
  <link rel="canonical" href="https://texly.chat/pricing">
  ```
- [ ] Use absolute URLs
- [ ] Add to all public pages

---

### Step 11: Internal Linking Strategy

#### 11.1 Link Structure
- [ ] From landing page:
  - [ ] Link to pricing, about, blog, use cases
- [ ] From blog posts:
  - [ ] Link to related posts
  - [ ] Link to pricing (CTA)
  - [ ] Link to sign-up
- [ ] From pricing:
  - [ ] Link to features
  - [ ] Link to FAQ
  - [ ] Link to sign-up

#### 11.2 Breadcrumbs
- [ ] Add breadcrumb navigation to blog:
  ```
  Home > Blog > [Post Title]
  ```
- [ ] Add structured data for breadcrumbs

---

## Testing Tasks

### SEO Audit
- [ ] Run SEO audit with tools:
  - [ ] Ahrefs (if budget allows)
  - [ ] SEMrush (free trial)
  - [ ] Screaming Frog (desktop tool)
- [ ] Check for:
  - [ ] Broken links (404s)
  - [ ] Missing meta descriptions
  - [ ] Duplicate titles
  - [ ] Slow pages (>3s load time)

### Lighthouse Testing
- [ ] Run Lighthouse on:
  - [ ] Landing page
  - [ ] Pricing page
  - [ ] Blog post (if exists)
- [ ] Fix all issues below 90 score

### Mobile Testing
- [ ] Test on real mobile devices
- [ ] Use Chrome DevTools mobile emulation
- [ ] Check Core Web Vitals on mobile

### Social Sharing Testing
- [ ] Test OG tags with:
  - [ ] Facebook Sharing Debugger
  - [ ] Twitter Card Validator
  - [ ] LinkedIn Post Inspector
- [ ] Ensure images display correctly

---

## Documentation

### Step 12: Update Documentation

#### 12.1 SEO Best Practices Guide
- [ ] Document SEO strategy for content team:
  - [ ] Title tag format: "[Page] - Texly"
  - [ ] Meta description length: 150-160 characters
  - [ ] Header hierarchy (H1 → H2 → H3)
  - [ ] Image alt text guidelines

#### 12.2 Content Guidelines
- [ ] Create content style guide:
  - [ ] Keyword research process
  - [ ] Target keywords per page
  - [ ] Internal linking strategy
  - [ ] Blog post checklist (meta, images, links, CTA)

---

## Success Metrics

- [ ] Google Search Console: 100+ indexed pages (if blog exists)
- [ ] Lighthouse score >90 on all pages
- [ ] Core Web Vitals: "Good" rating
- [ ] Organic search traffic: Track growth month-over-month
- [ ] Keyword rankings: Track target keywords (e.g., "AI chatbot builder")
- [ ] Backlinks: Acquire 10+ quality backlinks (outreach, guest posts)
- [ ] Social shares: Track shares of blog posts

---

## Long-Term SEO Tasks (Post-Phase 9)

- [ ] Content marketing: Publish 2-4 blog posts/month
- [ ] Link building: Outreach to relevant sites
- [ ] Guest posting: Write for AI/tech blogs
- [ ] Community engagement: Answer questions on Reddit, Quora, forums
- [ ] Press releases: Announce major features or milestones
- [ ] Partnerships: Collaborate with complementary SaaS tools
- [ ] Video content: YouTube tutorials (embeds on site for rich snippets)

---

**Created**: February 11, 2025
