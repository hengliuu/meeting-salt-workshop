'use client';

import Image from 'next/image';
import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

export default function LoginPage() {
  const [account, setAccount] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    
    // Simulate login logic
    console.log('Login attempt:', { account, password });
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    setIsLoading(false);
    // Handle success/error here
  };

  return (
    <div className="min-h-screen flex">
      {/* Left Side - Neon Art Background */}
      <div className="hidden lg:flex lg:w-1/2 relative overflow-hidden">
        {/* Animated Neon Background */}
        <div className="absolute inset-0 bg-neon-gradient">
          {/* Animated Neon Shapes */}
          <div className="absolute inset-0">
            {/* Large floating orbs */}
            <div className="absolute top-1/4 left-1/4 w-32 h-32 bg-gradient-to-r from-pink-500 to-violet-500 rounded-full blur-xl opacity-60 neon-pulse"></div>
            <div className="absolute top-3/4 right-1/4 w-40 h-40 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-full blur-xl opacity-50 neon-float"></div>
            <div className="absolute top-1/2 left-1/2 w-24 h-24 bg-gradient-to-r from-violet-500 to-purple-500 rounded-full blur-lg opacity-70 animate-ping"></div>
            
            {/* Geometric neon lines */}
            <div className="absolute top-0 left-0 w-full h-full">
              <svg className="w-full h-full opacity-30" viewBox="0 0 400 400">
                <defs>
                  <linearGradient id="neonGradient1" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style={{stopColor: '#ff00ff', stopOpacity: 1}} />
                    <stop offset="100%" style={{stopColor: '#00ffff', stopOpacity: 1}} />
                  </linearGradient>
                  <linearGradient id="neonGradient2" x1="0%" y1="100%" x2="100%" y2="0%">
                    <stop offset="0%" style={{stopColor: '#8b5cf6', stopOpacity: 1}} />
                    <stop offset="100%" style={{stopColor: '#06b6d4', stopOpacity: 1}} />
                  </linearGradient>
                </defs>
                <path 
                  d="M50 350 Q 200 50 350 350" 
                  stroke="url(#neonGradient1)" 
                  strokeWidth="3" 
                  fill="none"
                  className="animate-pulse"
                />
                <path 
                  d="M50 50 Q 200 350 350 50" 
                  stroke="url(#neonGradient2)" 
                  strokeWidth="2" 
                  fill="none"
                  className="animate-pulse"
                  style={{animationDelay: '1s'}}
                />
                <circle 
                  cx="100" 
                  cy="100" 
                  r="30" 
                  stroke="url(#neonGradient1)" 
                  strokeWidth="2" 
                  fill="none"
                  className="neon-rotate"
                />
                <circle 
                  cx="300" 
                  cy="300" 
                  r="40" 
                  stroke="url(#neonGradient2)" 
                  strokeWidth="2" 
                  fill="none"
                  className="neon-rotate"
                  style={{animationDirection: 'reverse', animationDelay: '0.5s'}}
                />
              </svg>
            </div>
            
            {/* Floating particles */}
            <div className="absolute top-20 left-20 w-2 h-2 bg-pink-400 rounded-full animate-bounce" style={{animationDelay: '0.5s'}}></div>
            <div className="absolute top-40 right-32 w-1 h-1 bg-cyan-400 rounded-full animate-ping" style={{animationDelay: '1s'}}></div>
            <div className="absolute bottom-32 left-40 w-3 h-3 bg-violet-400 rounded-full animate-pulse" style={{animationDelay: '1.5s'}}></div>
            <div className="absolute bottom-40 right-20 w-2 h-2 bg-blue-400 rounded-full animate-bounce" style={{animationDelay: '2s'}}></div>
          </div>
          
          {/* Overlay gradient for depth */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-black/10"></div>
        </div>
      </div>

      {/* Right Side - Login Form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-8 bg-white relative">
        {/* Mobile background for small screens */}
        <div className="absolute inset-0 lg:hidden bg-gradient-to-br from-fuchsia-600/10 via-purple-600/10 to-blue-600/10">
          {/* Subtle mobile neon effects */}
          <div className="absolute top-10 right-10 w-16 h-16 bg-gradient-to-r from-pink-400 to-violet-400 rounded-full blur-lg opacity-30 animate-pulse"></div>
          <div className="absolute bottom-20 left-10 w-12 h-12 bg-gradient-to-r from-blue-400 to-cyan-400 rounded-full blur-md opacity-40 neon-float"></div>
        </div>
        
        {/* SALT Logo - Top Left */}
        <div className="absolute top-6 left-6 z-10">
          <Image
            src="/logo-salt.svg"
            alt="SALT Logo"
            width={120}
            height={60}
            priority
            className="h-auto"
          />
        </div>

        {/* Login Form Card */}
        <div className="w-full max-w-md relative z-10">
          <div className="glass rounded-2xl shadow-2xl p-8 bg-white/80 backdrop-blur-sm border border-white/20">
            {/* Title */}
            <div className="text-center mb-8">
              <h1 className="text-3xl font-bold bg-gradient-to-r from-fuchsia-600 via-purple-600 to-blue-600 bg-clip-text text-transparent">
                Account Login
              </h1>
            </div>

            {/* Login Form */}
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Account Field */}
              <div>
                <label htmlFor="account" className="block text-sm font-medium text-gray-700 mb-2">
                  Account
                </label>
                <Input
                  id="account"
                  type="text"
                  value={account}
                  onChange={(e) => setAccount(e.target.value)}
                  placeholder="Enter your account"
                  required
                  aria-describedby="account-help"
                />
              </div>

              {/* Password Field */}
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                  Password
                </label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter your password"
                  required
                  aria-describedby="password-help"
                />
              </div>

              {/* Login Button */}
              <Button
                type="submit"
                variant="neon"
                size="lg"
                className="w-full"
                disabled={isLoading}
                aria-label="Submit login form"
              >
                {isLoading ? (
                  <div className="flex items-center gap-2">
                    <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    Logging in...
                  </div>
                ) : (
                  'Login'
                )}
              </Button>
            </form>

            {/* Additional Elements */}
            <div className="mt-6 text-center">
              <a 
                href="#" 
                className="text-sm text-gray-500 hover:text-purple-600 transition-colors duration-200"
              >
                Forgot your password?
              </a>
            </div>
          </div>

          {/* Decorative Elements */}
          <div className="absolute -top-4 -right-4 w-8 h-8 bg-gradient-to-r from-pink-500 to-violet-500 rounded-full blur-sm opacity-60 animate-pulse"></div>
          <div className="absolute -bottom-4 -left-4 w-6 h-6 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-full blur-sm opacity-50 animate-bounce"></div>
        </div>
      </div>
    </div>
  );
}
