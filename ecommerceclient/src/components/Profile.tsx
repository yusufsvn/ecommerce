import { Container, Typography, Button, Box, Paper, Avatar, Divider, Grid } from '@mui/material';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router';
import LocalShippingOutlinedIcon from '@mui/icons-material/LocalShippingOutlined';
import LocationOnOutlinedIcon from '@mui/icons-material/LocationOnOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';

const Profile = () => {
    const { userName, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/');
    };

    return (
        <Container component="main" maxWidth="md" sx={{ py: { xs: 4, md: 8 } }}>
            {/* Hero / Header Section */}
            <Box sx={{ mb: 6, display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center' }}>
                <Avatar 
                    sx={{ 
                        width: 100, height: 100, mb: 3, 
                        bgcolor: 'primary.main', fontSize: '3rem', fontWeight: 'bold',
                        boxShadow: '0 8px 24px rgba(84, 117, 150, 0.3)'
                    }}
                >
                    {userName ? userName.charAt(0).toUpperCase() : 'M'}
                </Avatar>
                <Typography component="h1" variant="h3" fontWeight="900" gutterBottom>
                    Hoş geldin, {userName || "Misafir"}!
                </Typography>
                <Typography variant="subtitle1" color="text.secondary" sx={{ maxWidth: 500 }}>
                    Hesap bilgilerinizi, siparişlerinizi ve adreslerinizi bu ekrandan kolayca yönetebilirsiniz.
                </Typography>
            </Box>

            <Divider sx={{ mb: 6 }} />

            {/* Quick Actions Grid */}
            <Grid container spacing={4} sx={{ mb: 6 }}>
                {/* Siparişlerim Kartı */}
                <Grid size={{ xs: 12, sm: 6 }}>
                    <Paper 
                        elevation={0}
                        onClick={() => navigate('/profile/orders')}
                        sx={{ 
                            p: 4, height: '100%', display: 'flex', flexDirection: 'column',
                            borderRadius: 4, border: '1px solid', borderColor: 'divider',
                            cursor: 'pointer', transition: 'all 0.3s ease',
                            '&:hover': { 
                                borderColor: 'primary.main',
                                boxShadow: '0 12px 30px rgba(84, 117, 150, 0.1)',
                                transform: 'translateY(-4px)'
                            }
                        }}
                    >
                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                            <Box sx={{ p: 1.5, borderRadius: 3, bgcolor: 'primary.50', color: 'primary.main', mr: 2 }}>
                                <LocalShippingOutlinedIcon sx={{ fontSize: 32 }} />
                            </Box>
                            <Typography variant="h5" fontWeight="bold">
                                Siparişlerim
                            </Typography>
                        </Box>
                        <Typography variant="body1" color="text.secondary" sx={{ flexGrow: 1, mb: 3 }}>
                            Geçmiş siparişlerinizi görüntüleyin, kargo durumunu takip edin ve detayları inceleyin.
                        </Typography>
                        <Typography variant="button" color="primary" fontWeight="bold" sx={{ display: 'flex', alignItems: 'center' }}>
                            Siparişlere Git &rarr;
                        </Typography>
                    </Paper>
                </Grid>

                {/* Adres Bilgilerim Kartı */}
                <Grid size={{ xs: 12, sm: 6 }}>
                    <Paper 
                        elevation={0}
                        onClick={() => navigate('/profile/address')}
                        sx={{ 
                            p: 4, height: '100%', display: 'flex', flexDirection: 'column',
                            borderRadius: 4, border: '1px solid', borderColor: 'divider',
                            cursor: 'pointer', transition: 'all 0.3s ease',
                            '&:hover': { 
                                borderColor: 'primary.main',
                                boxShadow: '0 12px 30px rgba(84, 117, 150, 0.1)',
                                transform: 'translateY(-4px)'
                            }
                        }}
                    >
                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                            <Box sx={{ p: 1.5, borderRadius: 3, bgcolor: 'primary.50', color: 'primary.main', mr: 2 }}>
                                <LocationOnOutlinedIcon sx={{ fontSize: 32 }} />
                            </Box>
                            <Typography variant="h5" fontWeight="bold">
                                Adres Bilgilerim
                            </Typography>
                        </Box>
                        <Typography variant="body1" color="text.secondary" sx={{ flexGrow: 1, mb: 3 }}>
                            Teslimat ve fatura adreslerinizi güncelleyin ve güvenli alışverişin keyfini çıkarın.
                        </Typography>
                        <Typography variant="button" color="primary" fontWeight="bold" sx={{ display: 'flex', alignItems: 'center' }}>
                            Adresleri Yönet &rarr;
                        </Typography>
                    </Paper>
                </Grid>
            </Grid>

            {/* Logout Section */}
            <Box sx={{ textAlign: 'center' }}>
                <Button
                    variant="text"
                    color="error"
                    onClick={handleLogout}
                    startIcon={<LogoutOutlinedIcon />}
                    sx={{ px: 4, py: 1.5, borderRadius: 8, fontWeight: 'bold' }}
                >
                    Güvenli Çıkış Yap
                </Button>
            </Box>
        </Container>
    );
};

export default Profile;
