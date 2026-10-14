import { Outlet } from 'react-router';
import Navbar from './Navbar';
import Footer from './Footer';
import { Box } from '@mui/material';

interface PublicLayoutProps {
    getTotalItems: () => number;
}

export default function PublicLayout({ getTotalItems }: PublicLayoutProps) {
    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            <Navbar getTotalItems={getTotalItems} />
            <Box sx={{ flexGrow: 1 }}>
                <Outlet />
            </Box>
            <Footer />
        </Box>
    );
}
