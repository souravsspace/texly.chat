# Phase 8: Marketing & Legal UI

## Goal
Build customer-facing marketing pages, improve dashboard UI/UX, and ensure legal compliance with privacy and terms pages.

---

## Frontend Tasks - Marketing Pages

### Step 1: Landing Page

#### 1.1 Hero Section
- [ ] Create `ui/src/routes/index.tsx` (public landing page):
  - [ ] **Hero Section**:
    - [ ] Compelling headline: "Build AI Chatbots Trained on Your Data"
    - [ ] Subheadline: "Create custom chatbots with RAG-powered responses in minutes. No coding required."
    - [ ] Primary CTA: "Start Free" → `/signup`
    - [ ] Secondary CTA: "View Demo" → Demo video or interactive widget
    - [ ] Hero image/animation (chatbot in action)
    - [ ] Trust indicators (e.g., "No credit card required")

#### 1.2 Features Section
- [ ] **Key Features Grid** (4-6 features):
  - [ ] **RAG-Powered AI**: "Accurate answers from your content"
  - [ ] **Multi-Source Support**: "Upload PDFs, Excel, URLs, text"
  - [ ] **Embeddable Widget**: "Add to any website with one line of code"
  - [ ] **Real-Time Chat**: "Streaming responses for instant feedback"
  - [ ] **Customizable**: "Match your brand colors and style"
  - [ ] **Analytics**: "Track conversations and improve over time"
  - [ ] Each feature with icon, title, description

#### 1.3 How It Works Section
- [ ] **Step-by-step Process** (3 steps):
  1. **Upload Your Data**: "Add URLs, files, or text to train your bot"
  2. **Customize Appearance**: "Match your brand with colors and settings"
  3. **Embed Anywhere**: "Copy one line of code and go live"
  - [ ] Visual illustrations for each step

#### 1.4 Social Proof
- [ ] **Testimonials** (if available):
  - [ ] Customer quotes
  - [ ] Company logos
  - [ ] Use case examples
- [ ] **Stats** (update as you grow):
  - [ ] "10,000+ chatbots created"
  - [ ] "1M+ conversations powered"
  - [ ] "99.9% uptime"

#### 1.5 Pricing Teaser
- [ ] **Pricing Preview**:
  - [ ] "Start Free, Upgrade Anytime"
  - [ ] Show Free and Pro tiers (simplified)
  - [ ] CTA: "View Full Pricing" → `/pricing`

#### 1.6 Final CTA
- [ ] **Bottom CTA Section**:
  - [ ] "Ready to Build Your Chatbot?"
  - [ ] CTA button: "Start Free Today"
  - [ ] Secondary link: "Book a Demo" (enterprise)

---

### Step 2: About Page

#### 2.1 About Page Content
- [ ] Create `ui/src/routes/about/index.tsx`:
  - [ ] **Mission Statement**:
    - [ ] Why Texly exists
    - [ ] Problem being solved
  - [ ] **Team** (if applicable):
    - [ ] Founder/team photos
    - [ ] Short bios
  - [ ] **Technology**:
    - [ ] Built with cutting-edge AI (OpenAI GPT)
    - [ ] RAG for accuracy
    - [ ] Privacy-first (data stays secure)
  - [ ] **Values**:
    - [ ] Transparency
    - [ ] User privacy
    - [ ] Innovation

---

### Step 3: Use Cases Page

#### 3.1 Use Cases Gallery
- [ ] Create `ui/src/routes/use-cases/index.tsx`:
  - [ ] **Customer Support**: "24/7 automated support with accurate answers"
  - [ ] **Documentation Helper**: "Let users search your docs with natural language"
  - [ ] **Lead Generation**: "Capture leads while answering questions"
  - [ ] **Internal Knowledge Base**: "Empower employees with instant access to company knowledge"
  - [ ] **E-commerce**: "Product recommendations and FAQs"
  - [ ] Each use case with:
    - [ ] Icon or illustration
    - [ ] Title
    - [ ] Description
    - [ ] "See Example" link or demo

---

### Step 4: Blog/Resources (Optional)

#### 4.1 Blog Setup
- [ ] Create `ui/src/routes/blog/index.tsx` (optional for SEO):
  - [ ] List of blog posts
  - [ ] Topics: AI, chatbots, RAG, tutorials
  - [ ] CMS integration (e.g., Markdown files or headless CMS)
- [ ] Individual blog post route: `ui/src/routes/blog/$slug.tsx`

---

## Frontend Tasks - Legal Pages

### Step 5: Terms of Service

#### 5.1 Terms Page
- [ ] Create `ui/src/routes/legal/terms.tsx`:
  - [ ] **Sections**:
    1. **Acceptance of Terms**: Who can use the service
    2. **User Accounts**: Registration requirements
    3. **Service Description**: What Texly provides
    4. **User Responsibilities**: Acceptable use policy
    5. **Payment Terms**: Billing, refunds, cancellations
    6. **Data Ownership**: User owns their data
    7. **Intellectual Property**: Texly owns the platform
    8. **Limitation of Liability**: Legal disclaimers
    9. **Termination**: Account closure policy
    10. **Governing Law**: Jurisdiction
    11. **Changes to Terms**: Right to update
    12. **Contact Information**: How to reach support
  - [ ] Use simple, readable language (avoid legalese where possible)
  - [ ] Add "Last Updated" date

#### 5.2 Terms Template
- [ ] Use open-source template (e.g., from Basecamp, GitHub)
- [ ] Customize for Texly.Chat specifics
- [ ] Review with legal counsel (recommended)

---

### Step 6: Privacy Policy

#### 6.1 Privacy Page
- [ ] Create `ui/src/routes/legal/privacy.tsx`:
  - [ ] **Sections**:
    1. **Information We Collect**:
       - [ ] Account info (email, name)
       - [ ] OAuth data (if using Google/GitHub)
       - [ ] Uploaded content (URLs, files, text)
       - [ ] Usage data (analytics, logs)
       - [ ] Cookies and tracking
    2. **How We Use Information**:
       - [ ] Provide the service
       - [ ] Improve features
       - [ ] Billing and payments
       - [ ] Security and fraud prevention
    3. **Data Sharing**:
       - [ ] We don't sell your data
       - [ ] Third-party services (OpenAI, MinIO, Polar)
       - [ ] Legal requirements (if compelled)
    4. **Data Storage & Security**:
       - [ ] Encrypted in transit (HTTPS)
       - [ ] Stored securely (SQLite, MinIO)
       - [ ] Data retention policy
    5. **Your Rights**:
       - [ ] Access your data
       - [ ] Delete your account (GDPR right to erasure)
       - [ ] Export your data
       - [ ] Opt-out of emails
    6. **Cookies**:
       - [ ] Essential cookies (authentication)
       - [ ] Analytics cookies (optional)
       - [ ] How to disable cookies
    7. **GDPR Compliance** (if EU users):
       - [ ] Legal basis for processing
       - [ ] Data controller info
       - [ ] DPO contact (if applicable)
    8. **CCPA Compliance** (if California users):
       - [ ] Right to know
       - [ ] Right to delete
       - [ ] Right to opt-out of sale (N/A since we don't sell)
    9. **Changes to Policy**: Right to update
    10. **Contact**: Privacy questions email
  - [ ] Add "Last Updated" date

#### 6.2 Privacy Template
- [ ] Use template from privacy policy generators (e.g., Termly, iubenda)
- [ ] Customize for Texly.Chat
- [ ] Ensure GDPR/CCPA compliance
- [ ] Review with legal counsel

---

### Step 7: Cookie Policy

#### 7.1 Cookie Policy Page
- [ ] Create `ui/src/routes/legal/cookies.tsx`:
  - [ ] **What Are Cookies**: Explanation
  - [ ] **Cookies We Use**:
    - [ ] Essential: Authentication tokens (JWT)
    - [ ] Functional: User preferences
    - [ ] Analytics: Usage tracking (if using Google Analytics, Plausible, etc.)
  - [ ] **How to Manage Cookies**: Browser settings
  - [ ] **Third-Party Cookies**: OAuth providers, analytics
  - [ ] Link to Privacy Policy

#### 7.2 Cookie Consent Banner
- [ ] Create `ui/src/components/cookie-banner.tsx`:
  - [ ] Show on first visit
  - [ ] "We use cookies to improve your experience"
  - [ ] Buttons: "Accept All" | "Reject Non-Essential" | "Customize"
  - [ ] Link to Cookie Policy
  - [ ] Store consent in localStorage
  - [ ] Conditionally load analytics scripts based on consent

---

### Step 8: Contact & Support Pages

#### 8.1 Contact Page
- [ ] Create `ui/src/routes/contact/index.tsx`:
  - [ ] Contact form:
    - [ ] Fields: Name, Email, Subject, Message
    - [ ] Dropdown: "Inquiry Type" (General, Support, Enterprise, Press)
    - [ ] Submit to `/api/contact/general`
    - [ ] Auto-reply confirmation email
  - [ ] Alternative contact methods:
    - [ ] Email: support@texly.chat
    - [ ] Twitter/X: @texlychat (if applicable)
  - [ ] FAQ link

#### 8.2 Support/FAQ Page
- [ ] Create `ui/src/routes/support/faq.tsx`:
  - [ ] Common questions:
    - [ ] "How do I create a bot?"
    - [ ] "What file formats are supported?"
    - [ ] "Can I embed on multiple websites?"
    - [ ] "How is billing calculated?"
    - [ ] "How do I cancel my subscription?"
    - [ ] "Is my data secure?"
  - [ ] Accordion-style answers (collapsible)
  - [ ] Search functionality (optional)
  - [ ] "Still need help?" → Contact form

---

## Frontend Tasks - Dashboard UI/UX Improvements

### Step 9: Dashboard Enhancements

#### 9.1 Onboarding Flow for New Users
- [ ] Create `ui/src/components/onboarding/onboarding-wizard.tsx`:
  - [ ] Triggered on first login
  - [ ] Step 1: "Welcome to Texly! Let's create your first bot."
  - [ ] Step 2: "Add a data source" (URL, file, or text)
  - [ ] Step 3: "Customize appearance"
  - [ ] Step 4: "Embed your bot" (show code snippet)
  - [ ] Step 5: "You're all set!" (link to dashboard)
  - [ ] Skippable with "Skip for now" button
  - [ ] Store completion status in user preferences

#### 9.2 Empty States
- [ ] Improve empty state UI:
  - [ ] **No Bots**: Illustration + "Create Your First Bot" CTA
  - [ ] **No Sources**: "Add a data source to get started"
  - [ ] **No Messages**: "Start a conversation to test your bot"
  - [ ] Use friendly copy and visuals

#### 9.3 Dashboard Home
- [ ] Update `ui/src/routes/dashboard/index.tsx`:
  - [ ] **Welcome message**: "Welcome back, [Name]!"
  - [ ] **Quick stats**:
    - [ ] Total bots
    - [ ] Messages this month
    - [ ] Active sources
    - [ ] Credits remaining (Pro tier)
  - [ ] **Recent activity**:
    - [ ] Latest chat messages
    - [ ] Recently updated bots
  - [ ] **Quick actions**:
    - [ ] Create New Bot
    - [ ] View Usage
    - [ ] Upgrade (if Free tier)
  - [ ] **Tips & Tutorials** (carousel or cards):
    - [ ] "How to optimize your bot for better answers"
    - [ ] "Best practices for data sources"

#### 9.4 Navigation Improvements
- [ ] Update `ui/src/components/layout/sidebar.tsx`:
  - [ ] Clear icons for each section
  - [ ] Active state highlighting
  - [ ] Tooltips on hover
  - [ ] Add sections:
    - [ ] Dashboard (home)
    - [ ] Bots
    - [ ] Usage & Billing
    - [ ] Settings
    - [ ] Support
  - [ ] Footer: Tier badge (Free/Pro/Enterprise)

#### 9.5 Bot Management UI
- [ ] Improve `ui/src/routes/dashboard/bots/index.tsx`:
  - [ ] Grid view with bot cards (name, sources count, message count)
  - [ ] List view option (toggle)
  - [ ] Search and filter bots
  - [ ] Sort by: Name, Date Created, Most Active
  - [ ] Batch actions: Delete, Archive

#### 9.6 Settings Page
- [ ] Update `ui/src/routes/dashboard/settings.tsx`:
  - [ ] Tabs: Profile, Security, Notifications, Billing, API Keys
  - [ ] **Profile**: Name, Email, Avatar, Connected OAuth
  - [ ] **Security**: Change password, 2FA (future)
  - [ ] **Notifications**: Email preferences (marketing, billing, support)
  - [ ] **Billing**: Redirect to billing/subscription page
  - [ ] **API Keys** (future): Generate API keys for integrations

---

## Backend Tasks

### Step 10: Contact Form Endpoints

#### 10.1 General Contact Endpoint
- [ ] Create `internal/handlers/contact/contact_handler.go`:
  - [ ] `POST /api/contact/general`:
    - [ ] Validate input
    - [ ] Store in database (optional)
    - [ ] Send email to support team
    - [ ] Auto-reply to user
    - [ ] Return success response

#### 10.2 Enterprise Inquiry Endpoint
- [ ] (Already covered in Phase 7)
- [ ] `POST /api/contact/enterprise`:
  - [ ] Capture enterprise-specific details
  - [ ] Send to sales email
  - [ ] Auto-reply with next steps

---

## Design & Styling

### Step 11: Design System Consistency

#### 11.1 Update Design Tokens
- [ ] Use TailwindCSS v4 variables for consistency:
  - [ ] Primary brand color
  - [ ] Secondary color
  - [ ] Accent color
  - [ ] Success, warning, error colors
  - [ ] Typography scale
  - [ ] Spacing scale

#### 11.2 Component Library
- [ ] Ensure all shadcn/ui components match brand:
  - [ ] Buttons (primary, secondary, ghost, destructive)
  - [ ] Forms (inputs, selects, checkboxes)
  - [ ] Cards
  - [ ] Modals/Dialogs
  - [ ] Toasts/Notifications

#### 11.3 Illustrations & Icons
- [ ] Add illustrations:
  - [ ] Hero section
  - [ ] Features section
  - [ ] Empty states
  - [ ] Onboarding steps
  - [ ] Error pages (404, 500)
- [ ] Icon library (Lucide React, Heroicons, or custom)

---

## SEO Preparation (Basic)

### Step 12: Meta Tags for Marketing Pages

#### 12.1 Add Meta Tags to Each Page
- [ ] Landing page meta:
  - [ ] `<title>Texly - Build AI Chatbots Trained on Your Data</title>`
  - [ ] `<meta name="description" content="...">`
- [ ] Pricing page meta
- [ ] About page meta
- [ ] Blog posts meta (if implemented)

#### 12.2 Open Graph Tags
- [ ] Add OG tags for social sharing:
  ```html
  <meta property="og:title" content="Texly - AI Chatbots">
  <meta property="og:description" content="...">
  <meta property="og:image" content="/og-image.png">
  <meta property="og:url" content="https://texly.chat">
  ```

---

## Testing Tasks

### Manual Testing
- [ ] Test all marketing pages on mobile
- [ ] Test legal page readability
- [ ] Test contact form submissions
- [ ] Test onboarding flow from start to finish
- [ ] Test dashboard navigation and quick actions

### Accessibility Testing
- [ ] Check keyboard navigation
- [ ] Test screen reader compatibility
- [ ] Ensure color contrast meets WCAG AA standards
- [ ] Add ARIA labels where needed

### Browser Testing
- [ ] Test on Chrome, Firefox, Safari, Edge
- [ ] Test on mobile browsers (iOS Safari, Chrome Mobile)

---

## Success Metrics

- [ ] Landing page conversion rate >5%
- [ ] Contact form submissions (measure inquiries)
- [ ] Onboarding completion rate >80%
- [ ] Dashboard engagement (users visit >3 pages/session)
- [ ] Legal compliance (no GDPR/CCPA complaints)
- [ ] Mobile traffic >30% (responsive design working)

---

## Documentation

- [ ] Update README with links to marketing pages
- [ ] Document design system in Storybook (optional)
- [ ] Create brand guidelines (logo usage, colors, typography)

---

**Created**: February 11, 2025
