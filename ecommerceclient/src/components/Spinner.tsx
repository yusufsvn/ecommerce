import React from 'react';
import { Box, CircularProgress, Typography } from '@mui/material';

interface SpinnerProps {
    message?: string;
}

const Spinner: React.FC<SpinnerProps> = ({ message = 'Yükleniyor...' }) => {
    return (
        <Box 
            sx={{ 
                display: 'flex', 
                flexDirection: 'column', 
                alignItems: 'center', 
                justifyContent: 'center', 
                minHeight: '40vh',
                py: 6
            }}
        >
            <CircularProgress size={60} thickness={4} color="primary" />
            <Typography variant="h6" sx={{ mt: 2, color: 'text.secondary', fontWeight: 500 }}>
                {message}
            </Typography>
        </Box>
    );
};

export default Spinner;
