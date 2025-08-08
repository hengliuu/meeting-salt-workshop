import Link from "next/link";
import Image from "next/image";

export default function Home() {
  return (
    <div className="font-sans grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 bg-gradient-to-br from-fuchsia-50 via-purple-50 to-blue-50">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        <Image
          src="/logo-salt.svg"
          alt="SALT Logo"
          width={200}
          height={100}
          priority
          className="h-auto"
        />
        
        <div className="text-center sm:text-left">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-fuchsia-600 via-purple-600 to-blue-600 bg-clip-text text-transparent mb-4">
            Welcome to SALT
          </h1>
          <p className="text-gray-600 text-lg mb-8">
            Secure Authentication & Login Technology
          </p>
        </div>

        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <Link
            href="/dashboard"
            className="rounded-full bg-gradient-to-r from-fuchsia-600 via-purple-600 to-blue-600 hover:from-fuchsia-700 hover:via-purple-700 hover:to-blue-700 text-white font-semibold text-sm sm:text-base h-12 px-8 flex items-center justify-center transition-all duration-300 transform hover:scale-105 hover:shadow-lg"
          >
            Go to Dashboard
          </Link>
          <Link
            href="/login"
            className="rounded-full border border-solid border-purple-200 hover:border-purple-300 transition-colors flex items-center justify-center hover:bg-purple-50 font-medium text-sm sm:text-base h-12 px-6 text-purple-700"
          >
            Login Page
          </Link>
        </div>
      </main>
      
      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center text-gray-500">
        <span className="text-sm">
          Powered by Next.js & TailwindCSS
        </span>
      </footer>
    </div>
  );
}
