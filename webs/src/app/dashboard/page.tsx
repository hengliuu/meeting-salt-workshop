'use client';

import Image from 'next/image';
import Link from 'next/link';
import { Button } from '@/components/ui/button';

export default function DashboardPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white">
      {/* Header Navigation */}
      <header className="bg-gradient-to-r from-fuchsia-600 via-purple-600 to-blue-600 text-white p-4">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Image
              src="/logo-salt.svg"
              alt="SALT Logo"
              width={80}
              height={40}
              className="h-auto"
            />
            <h1 className="text-xl font-bold">MEETING ROOM BOOKING SYSTEM</h1>
          </div>
          
          <nav className="hidden md:flex items-center gap-6">
            <div className="flex items-center gap-2">
              <input
                type="date"
                defaultValue="2025-08-08"
                className="px-3 py-1 rounded bg-white/20 border border-white/30 text-white placeholder-white/70"
              />
              <Button variant="outline" size="sm" className="text-white border-white/30 hover:bg-white/10">
                Go to
              </Button>
            </div>
            <Link href="#" className="hover:text-white/80 transition-colors">HELP</Link>
            <Link href="#" className="hover:text-white/80 transition-colors">ROOMS</Link>
            <Link href="#" className="hover:text-white/80 transition-colors">REPORT</Link>
            <div className="flex items-center gap-2">
              <span className="text-sm">Search:</span>
              <input
                type="text"
                className="px-2 py-1 rounded bg-white/20 border border-white/30 text-white placeholder-white/70 w-32"
                placeholder="Search..."
              />
            </div>
            <div className="flex items-center gap-2 text-sm">
              <span>joshua.salty</span>
              <Link href="/login">
                <Button variant="outline" size="sm" className="text-white border-white/30 hover:bg-white/10">
                  Log off
                </Button>
              </Link>
              <Link href="/">
                <Button variant="outline" size="sm" className="text-white border-white/30 hover:bg-white/10">
                  USER LIST
                </Button>
              </Link>
            </div>
          </nav>
        </div>
      </header>

      <div className="max-w-7xl mx-auto p-6 space-y-8">
        {/* Banner Section */}
        <section className="relative overflow-hidden rounded-3xl bg-gradient-to-r from-purple-900 via-blue-900 to-indigo-900 text-white p-8 md:p-12 fade-in-up">
          {/* Futuristic Background Graphics */}
          <div className="absolute inset-0 opacity-20 tech-grid">
            <div className="absolute top-10 left-10 w-32 h-32 border border-cyan-400 rounded-full animate-spin" style={{animationDuration: '20s'}}></div>
            <div className="absolute top-20 right-20 w-24 h-24 bg-gradient-to-r from-pink-500 to-violet-500 rounded-full blur-xl animate-pulse"></div>
            <div className="absolute bottom-10 left-1/4 w-40 h-40 border-2 border-blue-400 rotate-45 animate-bounce" style={{animationDuration: '3s'}}></div>
            <div className="absolute bottom-20 right-10 w-28 h-28 bg-gradient-to-r from-cyan-400 to-blue-500 rounded-full blur-lg opacity-60 animate-pulse" style={{animationDelay: '1s'}}></div>
            
            {/* Tech Grid Pattern */}
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-cyan-500/10 to-transparent transform skew-x-12"></div>
            <div className="absolute inset-0 bg-gradient-to-b from-transparent via-purple-500/10 to-transparent transform -skew-y-12"></div>
          </div>
          
          <div className="relative z-10 text-center space-y-6">
            <h1 className="text-5xl md:text-7xl font-bold bg-gradient-to-r from-cyan-400 via-blue-400 to-purple-400 bg-clip-text text-transparent">
              LMS Program Launching
            </h1>
            <p className="text-xl md:text-2xl text-gray-200 max-w-3xl mx-auto">
              Empowering Our Team's Future through:
            </p>
            <div className="flex flex-wrap justify-center gap-4 mt-8">
              <span className="px-6 py-3 bg-white/10 backdrop-blur-sm rounded-full border border-white/20 text-cyan-300">Innovation</span>
              <span className="px-6 py-3 bg-white/10 backdrop-blur-sm rounded-full border border-white/20 text-purple-300">Learning</span>
              <span className="px-6 py-3 bg-white/10 backdrop-blur-sm rounded-full border border-white/20 text-blue-300">Growth</span>
              <span className="px-6 py-3 bg-white/10 backdrop-blur-sm rounded-full border border-white/20 text-indigo-300">Excellence</span>
            </div>
          </div>
        </section>

        {/* Quote Section */}
        <section className="text-center py-12 bg-gradient-to-r from-gray-50 to-gray-100 rounded-2xl fade-in-up delay-1">
          <blockquote className="text-2xl md:text-4xl font-bold text-gray-800 mb-4">
            "Embrace challenges, Create solutions"
          </blockquote>
          <cite className="text-lg text-purple-600 font-semibold">
            – Marco Widjojo, CEO SALT
          </cite>
        </section>

        {/* Vision and Culture Cards */}
        <div className="grid md:grid-cols-2 gap-8 fade-in-up delay-2">
          {/* Vision Card */}
          <div className="bg-gradient-to-br from-purple-100 to-blue-100 p-8 rounded-2xl border border-purple-200 shadow-lg card-hover">
            <h3 className="text-2xl font-bold text-purple-800 mb-4 flex items-center gap-2">
              <div className="w-8 h-8 bg-purple-600 rounded-full flex items-center justify-center">
                <span className="text-white text-sm">🎯</span>
              </div>
              Vision
            </h3>
            <p className="text-lg text-gray-700 leading-relaxed">
              "To be the leading and most innovative technology company"
            </p>
          </div>

          {/* Company Culture Card */}
          <div className="bg-gradient-to-br from-cyan-100 to-blue-100 p-8 rounded-2xl border border-cyan-200 shadow-lg card-hover">
            <h3 className="text-2xl font-bold text-cyan-800 mb-6 flex items-center gap-2">
              <div className="w-8 h-8 bg-cyan-600 rounded-full flex items-center justify-center">
                <span className="text-white text-sm">🏢</span>
              </div>
              Company Culture
            </h3>
            <ul className="space-y-3">
              <li className="flex items-center gap-3 text-gray-700">
                <div className="w-2 h-2 bg-cyan-500 rounded-full"></div>
                Client Orientation
              </li>
              <li className="flex items-center gap-3 text-gray-700">
                <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                Innovation & Technology
              </li>
              <li className="flex items-center gap-3 text-gray-700">
                <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                Teamwork & Collaboration
              </li>
              <li className="flex items-center gap-3 text-gray-700">
                <div className="w-2 h-2 bg-indigo-500 rounded-full"></div>
                Open Communication
              </li>
              <li className="flex items-center gap-3 text-gray-700">
                <div className="w-2 h-2 bg-pink-500 rounded-full"></div>
                Social Responsibility
              </li>
            </ul>
          </div>
        </div>

        {/* News Section */}
        <section className="bg-white p-8 rounded-2xl shadow-lg border border-gray-200 fade-in-up delay-3 card-hover">
          <h3 className="text-2xl font-bold text-gray-800 mb-6 flex items-center gap-2">
            <div className="w-8 h-8 bg-orange-600 rounded-full flex items-center justify-center">
              <span className="text-white text-sm">📰</span>
            </div>
            Latest News
          </h3>
          <div className="grid md:grid-cols-2 gap-6">
            <article className="p-6 bg-gradient-to-r from-orange-50 to-red-50 rounded-xl border border-orange-200 hover:shadow-md transition-shadow">
              <h4 className="text-lg font-semibold text-gray-800 mb-2">
                The Ins & Outs of IT Solutions: Services and Its Examples
              </h4>
              <p className="text-gray-600 text-sm mb-3">Jul 2024</p>
              <p className="text-gray-700 text-sm">
                Explore comprehensive IT solutions and understand how they can transform your business operations...
              </p>
            </article>
            <article className="p-6 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl border border-blue-200 hover:shadow-md transition-shadow">
              <h4 className="text-lg font-semibold text-gray-800 mb-2">
                IT Outsourcing and Its Prospects
              </h4>
              <p className="text-gray-600 text-sm mb-3">24 Jun 2024</p>
              <p className="text-gray-700 text-sm">
                Understanding the future of IT outsourcing and how it benefits modern businesses...
              </p>
            </article>
          </div>
        </section>

        {/* Quick Links and Employee Benefits */}
        <div className="grid md:grid-cols-2 gap-8 fade-in-up delay-4">
          {/* Quick Links */}
          <div className="bg-gradient-to-br from-green-100 to-emerald-100 p-8 rounded-2xl border border-green-200 shadow-lg card-hover">
            <h3 className="text-2xl font-bold text-green-800 mb-6 flex items-center gap-2">
              <div className="w-8 h-8 bg-green-600 rounded-full flex items-center justify-center">
                <span className="text-white text-sm">🔗</span>
              </div>
              Quick Links
            </h3>
            <div className="space-y-4">
              <Link href="#" className="flex items-center gap-3 p-3 bg-white/50 rounded-lg hover:bg-white/70 transition-colors">
                <div className="w-10 h-10 bg-blue-500 rounded-lg flex items-center justify-center">
                  <span className="text-white text-sm">📄</span>
                </div>
                <span className="font-medium text-gray-700">Documents</span>
              </Link>
              <Link href="#" className="flex items-center gap-3 p-3 bg-white/50 rounded-lg hover:bg-white/70 transition-colors">
                <div className="w-10 h-10 bg-purple-500 rounded-lg flex items-center justify-center">
                  <span className="text-white text-sm">📸</span>
                </div>
                <span className="font-medium text-gray-700">Photos</span>
              </Link>
              <Link href="#" className="flex items-center gap-3 p-3 bg-white/50 rounded-lg hover:bg-white/70 transition-colors">
                <div className="w-10 h-10 bg-indigo-500 rounded-lg flex items-center justify-center">
                  <span className="text-white text-sm">⚙️</span>
                </div>
                <span className="font-medium text-gray-700">Others</span>
              </Link>
            </div>
          </div>

          {/* Employee Benefits */}
          <div className="bg-gradient-to-br from-yellow-100 to-orange-100 p-8 rounded-2xl border border-yellow-200 shadow-lg card-hover">
            <h3 className="text-2xl font-bold text-orange-800 mb-6 flex items-center gap-2">
              <div className="w-8 h-8 bg-orange-600 rounded-full flex items-center justify-center">
                <span className="text-white text-sm">🎁</span>
              </div>
              Employee Benefits
            </h3>
            <div className="space-y-4">
              <div className="flex items-center gap-3 p-4 bg-white/50 rounded-lg">
                <div className="w-12 h-12 bg-blue-500 rounded-lg flex items-center justify-center">
                  <span className="text-white text-lg">🏨</span>
                </div>
                <div>
                  <h4 className="font-semibold text-gray-700">Hotel</h4>
                  <p className="text-sm text-gray-600">Accommodation benefits</p>
                </div>
              </div>
              <div className="flex items-center gap-3 p-4 bg-white/50 rounded-lg">
                <div className="w-12 h-12 bg-green-500 rounded-lg flex items-center justify-center">
                  <span className="text-white text-lg">🏖️</span>
                </div>
                <div>
                  <h4 className="font-semibold text-gray-700">Resort</h4>
                  <p className="text-sm text-gray-600">Vacation packages</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Feedback and Wellness Section */}
        <div className="grid md:grid-cols-2 gap-8 fade-in-up">
          {/* Feedback Card */}
          <div className="bg-gradient-to-br from-pink-100 to-rose-100 p-8 rounded-2xl border border-pink-200 shadow-lg card-hover">
            <h3 className="text-2xl font-bold text-pink-800 mb-4 flex items-center gap-2">
              <div className="w-8 h-8 bg-pink-600 rounded-full flex items-center justify-center">
                <span className="text-white text-sm">💬</span>
              </div>
              Feedback
            </h3>
            <p className="text-gray-700 mb-6">
              "Please share your thoughts with us!"
            </p>
            <Button 
              variant="neon" 
              className="w-full bg-gradient-to-r from-pink-500 to-rose-500 hover:from-pink-600 hover:to-rose-600"
            >
              Send Feedback
            </Button>
          </div>

          {/* Wellness Promo */}
          <div className="bg-gradient-to-br from-cyan-100 to-teal-100 p-8 rounded-2xl border border-cyan-200 shadow-lg card-hover">
            <h3 className="text-2xl font-bold text-teal-800 mb-4 flex items-center gap-2">
              <div className="w-8 h-8 bg-teal-600 rounded-full flex items-center justify-center">
                <span className="text-white text-sm">🏃</span>
              </div>
              Bergerak Untuk Hidup Sehat
            </h3>
            
            {/* Sports/Athletes Imagery Placeholder */}
            <div className="mb-6 h-32 bg-gradient-to-r from-teal-200 to-cyan-200 rounded-lg flex items-center justify-center">
              <div className="text-6xl">🏃‍♂️🏊‍♀️🚴‍♂️</div>
            </div>
            
            <Button 
              variant="outline" 
              className="w-full border-teal-500 text-teal-700 hover:bg-teal-50"
            >
              See More
            </Button>
          </div>
        </div>

        {/* Footer */}
        <footer className="text-center py-8 text-gray-500">
          <p>&copy; 2025 SALT - Secure Authentication & Login Technology. All rights reserved.</p>
        </footer>
      </div>
    </div>
  );
}
