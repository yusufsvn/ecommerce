import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App'
import { BrowserRouter } from 'react-router'
import { AuthProvider } from './context/AuthContext'
import { ToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import GlobalLoader from './components/GlobalLoader'
import { GoogleOAuthProvider } from '@react-oauth/google'

const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <AuthProvider>
      <GoogleOAuthProvider clientId={clientId}>
        <App />
        <GlobalLoader />
        <ToastContainer position="top-right" autoClose={3000} />
      </GoogleOAuthProvider>
    </AuthProvider>
  </BrowserRouter>
)
