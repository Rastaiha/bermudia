import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route, Navigate, Link } from 'react-router-dom'
import TerritoryEditor from './components/TerritoryEditor'
import './index.css'

// Placeholder components for future admin panel pages
const Dashboard = () => (
  <div className="p-6">
    <h1 className="text-2xl font-bold text-gray-800 mb-4">Admin Dashboard</h1>
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-lg font-semibold mb-2">Territories</h2>
        <p className="text-gray-600">Manage game territories and islands</p>
        <Link 
          to="/territories" 
          className="inline-block mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Manage Territories
        </Link>
      </div>
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-lg font-semibold mb-2">Players</h2>
        <p className="text-gray-600">View and manage player accounts</p>
        <button className="inline-block mt-4 px-4 py-2 bg-gray-400 text-white rounded cursor-not-allowed">
          Coming Soon
        </button>
      </div>
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-lg font-semibold mb-2">Analytics</h2>
        <p className="text-gray-600">Game statistics and reports</p>
        <button className="inline-block mt-4 px-4 py-2 bg-gray-400 text-white rounded cursor-not-allowed">
          Coming Soon
        </button>
      </div>
    </div>
  </div>
)

const Layout = ({ children }: { children: React.ReactNode }) => (
  <div className="min-h-screen bg-gray-100">
    {/* Navigation Header */}
    <nav className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link to="/" className="text-xl font-bold text-gray-800">
              Education Game Admin
            </Link>
          </div>
          <div className="flex items-center space-x-4">
            <Link 
              to="/" 
              className="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
            >
              Dashboard
            </Link>
            <Link 
              to="/territories" 
              className="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
            >
              Territories
            </Link>
          </div>
        </div>
      </div>
    </nav>
    
    {/* Main Content */}
    <main className="flex-1">
      {children}
    </main>
  </div>
)

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout><Dashboard /></Layout>} />
        <Route path="/territories" element={<TerritoryEditor />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)