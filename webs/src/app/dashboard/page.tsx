'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';

const months = [
  { name: 'July 2025', days: 31, startDay: 2 }, // July 1st is Tuesday (index 2)
  { name: 'August 2025', days: 31, startDay: 5 }, // August 1st is Friday (index 5)
  { name: 'September 2025', days: 30, startDay: 1 }, // September 1st is Monday (index 1)
];

const timeSlots = [
  '07:00', '07:30', '08:00', '08:30', '09:00', '09:30', '10:00', '10:30', '11:00'
];

export default function DashboardPage() {
  const [selectedFloor, setSelectedFloor] = useState('1st Floor');
  
  const floors = ['1st Floor', '2nd Floor', '3rd Floor', '4th Floor'];
  
  const generateCalendar = (month: { name: string; days: number; startDay: number }) => {
    const calendar = [];
    const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    
    // Add header
    calendar.push(
      <div key={`${month.name}-header`} className="grid grid-cols-7 gap-1 mb-2">
        {daysOfWeek.map(day => (
          <div key={day} className="text-center text-xs font-semibold text-gray-600 p-1">
            {day}
          </div>
        ))}
      </div>
    );
    
    // Add calendar days
    const weeks = [];
    let currentWeek = [];
    
    // Add empty cells for days before month starts
    for (let i = 0; i < month.startDay; i++) {
      currentWeek.push(<div key={`empty-${i}`} className="p-1"></div>);
    }
    
    // Add month days
    for (let day = 1; day <= month.days; day++) {
      const isSelected = month.name === 'August 2025' && day === 8;
      const isToday = month.name === 'August 2025' && day === 15;
      
      currentWeek.push(
        <div
          key={day}
          className={`text-center p-1 text-xs cursor-pointer hover:bg-pink-100 rounded ${
            isSelected ? 'bg-pink-500 text-white font-bold' : 
            isToday ? 'bg-pink-200 text-pink-800 font-semibold' :
            'text-gray-700'
          }`}
        >
          {day}
        </div>
      );
      
      if (currentWeek.length === 7) {
        weeks.push(
          <div key={`week-${weeks.length}`} className="grid grid-cols-7 gap-1">
            {currentWeek}
          </div>
        );
        currentWeek = [];
      }
    }
    
    // Add remaining empty cells
    while (currentWeek.length < 7) {
      currentWeek.push(<div key={`empty-end-${currentWeek.length}`} className="p-1"></div>);
    }
    
    if (currentWeek.length > 0) {
      weeks.push(
        <div key={`week-${weeks.length}`} className="grid grid-cols-7 gap-1">
          {currentWeek}
        </div>
      );
    }
    
    calendar.push(...weeks);
    return calendar;
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header Navigation */}
      <header className="bg-pink-600 text-white p-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-6">
            <h1 className="text-lg font-bold">MEETING ROOM BOOKING SYSTEM</h1>
          </div>
          
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <input
                type="date"
                defaultValue="2025-08-08"
                className="px-2 py-1 rounded text-black text-sm"
              />
              <Button size="sm" className="bg-white text-pink-600 hover:bg-gray-100 text-sm px-3 py-1">
                Go to
              </Button>
            </div>
            <Link href="#" className="hover:text-pink-200 transition-colors text-sm">HELP</Link>
            <Link href="#" className="hover:text-pink-200 transition-colors text-sm">ROOMS</Link>
            <Link href="#" className="hover:text-pink-200 transition-colors text-sm">REPORT</Link>
            <div className="flex items-center gap-2">
              <span className="text-sm">Search:</span>
              <input
                type="text"
                className="px-2 py-1 rounded text-black text-sm w-24"
                placeholder=""
              />
            </div>
            <div className="flex items-center gap-2 text-sm">
              <span>joshua.salty</span>
              <Link href="/login">
                <Button size="sm" className="bg-white text-pink-600 hover:bg-gray-100 text-xs px-2 py-1">
                  Log off
                </Button>
              </Link>
              <Link href="/">
                <Button size="sm" className="bg-white text-pink-600 hover:bg-gray-100 text-xs px-2 py-1">
                  USER LIST
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </header>

      <div className="flex">
        {/* Left Sidebar */}
        <aside className="w-48 bg-white border-r border-gray-300">
          <div className="p-4">
            <h3 className="font-semibold text-gray-700 mb-4">Areas</h3>
            <ul className="space-y-2">
              {floors.map((floor) => (
                <li key={floor}>
                  <button
                    onClick={() => setSelectedFloor(floor)}
                    className={`text-left w-full p-2 rounded text-sm ${
                      selectedFloor === floor
                        ? 'text-pink-600 font-semibold bg-pink-50'
                        : 'text-gray-700 hover:text-pink-600 hover:bg-gray-50'
                    }`}
                  >
                    {floor}
                  </button>
                </li>
              ))}
            </ul>
          </div>
        </aside>

        {/* Main Content */}
        <main className="flex-1 p-4">
          {/* Calendar Section */}
          <div className="bg-white rounded border border-gray-300 mb-6">
            <div className="grid grid-cols-3 gap-4 p-4">
              {months.map((month) => (
                <div key={month.name} className="text-center">
                  <h4 className="font-semibold text-gray-800 mb-2 text-sm">{month.name}</h4>
                  {generateCalendar(month)}
                </div>
              ))}
            </div>
          </div>

          {/* Daily Schedule Section */}
          <div className="bg-white rounded border border-gray-300">
            {/* Schedule Header */}
            <div className="border-b border-gray-300 p-4">
              <div className="flex items-center justify-between">
                <h2 className="text-xl font-bold text-gray-800">Friday 08 August 2025</h2>
                <div className="flex items-center gap-4">
                  <Link href="#" className="text-pink-600 hover:text-pink-800 text-sm">« Go To Day Before</Link>
                  <Link href="#" className="text-pink-600 hover:text-pink-800 text-sm">Go To Today</Link>
                  <Link href="#" className="text-pink-600 hover:text-pink-800 text-sm">Go To Day After »</Link>
                </div>
              </div>
            </div>

            {/* Schedule Grid */}
            <div className="grid grid-cols-3 h-96">
              {/* Time Column */}
              <div className="border-r border-gray-300">
                <div className="bg-pink-600 text-white text-center py-2 font-semibold text-sm">
                  Time
                </div>
                {timeSlots.map((time) => (
                  <div key={time} className="border-b border-gray-200 p-2 text-sm text-gray-700 h-10 flex items-center">
                    {time}
                  </div>
                ))}
              </div>

              {/* Smart Solution Column */}
              <div className="border-r border-gray-300 relative">
                <div className="bg-pink-600 text-white text-center py-2 font-semibold text-sm">
                  Smart Solution (8)
                </div>
                <div className="relative">
                  {timeSlots.map((time, index) => (
                    <div key={time} className="border-b border-gray-200 h-10"></div>
                  ))}
                  
                  {/* Meeting Block - ICT Team Operational Agreement */}
                  <div 
                    className="absolute left-1 right-1 bg-purple-600 text-white text-xs p-1 rounded shadow-sm flex items-center justify-center"
                    style={{
                      top: '120px', // Starting around 09:00
                      height: '80px' // Spanning about 2 hours
                    }}
                  >
                    <div className="text-center">
                      <div className="font-medium">ICT Team Operational Agreement</div>
                      <div className="text-xs opacity-90 mt-1">🔄</div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Technology Column */}
              <div>
                <div className="bg-pink-600 text-white text-center py-2 font-semibold text-sm">
                  Technology (4)
                </div>
                {timeSlots.map((time) => (
                  <div key={time} className="border-b border-gray-200 h-10"></div>
                ))}
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
