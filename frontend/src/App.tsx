import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout/Layout'
import Home from './pages/Home'
import Confirm from './pages/Confirm'
import Progress from './pages/Progress'
import Report from './pages/Report'
import History from './pages/History'
import { ToastProvider } from './hooks/useToast'
import ErrorBoundary from './components/common/ErrorBoundary'

function App() {
  return (
    <ErrorBoundary>
      <ToastProvider>
        <BrowserRouter>
          <Layout>
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/confirm" element={<Confirm />} />
              <Route path="/progress/:id" element={<Progress />} />
              <Route path="/report/:id" element={<Report />} />
              <Route path="/history" element={<History />} />
            </Routes>
          </Layout>
        </BrowserRouter>
      </ToastProvider>
    </ErrorBoundary>
  )
}

export default App
