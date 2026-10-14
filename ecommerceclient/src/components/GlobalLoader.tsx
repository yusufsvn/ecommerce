import React, { useEffect, useState } from 'react';
import Backdrop from '@mui/material/Backdrop';
import CircularProgress from '@mui/material/CircularProgress';

const GlobalLoader: React.FC = () => {
    const [open, setOpen] = useState(false);
    // Keep a counter to handle concurrent API requests.
    const [requestCount, setRequestCount] = useState(0);

    useEffect(() => {
        const handleStart = () => {
            setRequestCount((prev) => prev + 1);
        };

        const handleEnd = () => {
            setRequestCount((prev) => Math.max(0, prev - 1));
        };

        window.addEventListener('api-request-start', handleStart);
        window.addEventListener('api-request-end', handleEnd);

        return () => {
            window.removeEventListener('api-request-start', handleStart);
            window.removeEventListener('api-request-end', handleEnd);
        };
    }, []);

    useEffect(() => {
        setOpen(requestCount > 0);
    }, [requestCount]);

    return (
        <Backdrop
            sx={(theme) => ({ color: '#fff', zIndex: theme.zIndex.drawer + 9999 })}
            open={open}
        >
            <CircularProgress color="inherit" />
        </Backdrop>
    );
};

export default GlobalLoader;
