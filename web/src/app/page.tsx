import Link from 'next/link';

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-white to-gray-50">
      {/* Header */}
      <header className="container mx-auto px-4 py-6">
        <nav className="flex items-center justify-between">
          <div className="text-2xl font-bold text-primary">HireMe</div>
          <div className="flex items-center gap-4">
            <Link 
              href="/dashboard" 
              className="text-sm text-muted-foreground hover:text-foreground"
            >
              Dashboard
            </Link>
            <Link
              href="/editor"
              className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              Create CV
            </Link>
          </div>
        </nav>
      </header>

      {/* Hero */}
      <section className="container mx-auto px-4 py-24 text-center">
        <h1 className="text-5xl font-bold tracking-tight text-foreground sm:text-6xl">
          Build Your Perfect CV
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg text-muted-foreground">
          Create stunning, professional CVs in minutes with our intuitive drag-and-drop 
          editor. Export to PDF, DOCX, or share directly with employers.
        </p>
        <div className="mt-10 flex items-center justify-center gap-4">
          <Link
            href="/editor"
            className="rounded-md bg-primary px-6 py-3 text-sm font-semibold text-primary-foreground shadow-sm hover:bg-primary/90"
          >
            Start Building — It&apos;s Free
          </Link>
          <Link
            href="#features"
            className="rounded-md px-6 py-3 text-sm font-semibold text-foreground hover:bg-muted"
          >
            Learn More →
          </Link>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="container mx-auto px-4 py-24">
        <h2 className="text-center text-3xl font-bold text-foreground">
          Everything You Need
        </h2>
        <div className="mt-12 grid gap-8 md:grid-cols-3">
          {/* Feature 1 */}
          <div className="rounded-lg border bg-card p-6">
            <div className="mb-4 text-3xl">🎨</div>
            <h3 className="text-xl font-semibold">3 Professional Templates</h3>
            <p className="mt-2 text-muted-foreground">
              Choose from Classic, Modern, or Minimal designs. Each template is 
              optimized for ATS systems.
            </p>
          </div>

          {/* Feature 2 */}
          <div className="rounded-lg border bg-card p-6">
            <div className="mb-4 text-3xl">🖱️</div>
            <h3 className="text-xl font-semibold">Drag & Drop Editor</h3>
            <p className="mt-2 text-muted-foreground">
              Easily rearrange sections, add custom content, and see changes 
              in real-time with our live preview.
            </p>
          </div>

          {/* Feature 3 */}
          <div className="rounded-lg border bg-card p-6">
            <div className="mb-4 text-3xl">📄</div>
            <h3 className="text-xl font-semibold">Multi-Format Export</h3>
            <p className="mt-2 text-muted-foreground">
              Export your CV as PDF, DOCX, or raw JSON. Perfect for any 
              application requirement.
            </p>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="container mx-auto px-4 py-24 text-center">
        <div className="rounded-2xl bg-primary/5 p-12">
          <h2 className="text-3xl font-bold text-foreground">
            Ready to Land Your Dream Job?
          </h2>
          <p className="mx-auto mt-4 max-w-xl text-muted-foreground">
            Join thousands of professionals who have already created their 
            perfect CV with HireMe.
          </p>
          <Link
            href="/editor"
            className="mt-8 inline-block rounded-md bg-primary px-8 py-3 text-sm font-semibold text-primary-foreground shadow-sm hover:bg-primary/90"
          >
            Get Started Now
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t bg-muted/50">
        <div className="container mx-auto px-4 py-8">
          <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
            <div className="text-sm text-muted-foreground">
              © {new Date().getFullYear()} HireMe. Open source under MIT license.
            </div>
            <div className="flex gap-6">
              <Link href="/privacy" className="text-sm text-muted-foreground hover:text-foreground">
                Privacy
              </Link>
              <Link href="/terms" className="text-sm text-muted-foreground hover:text-foreground">
                Terms
              </Link>
              <a 
                href="https://github.com/yourusername/hireme" 
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                GitHub
              </a>
            </div>
          </div>
        </div>
      </footer>
    </main>
  );
}
