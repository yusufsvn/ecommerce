import React, { useState } from 'react';
import { Container, TextField, Button, Typography, Box, Paper, Avatar } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import { useAuth } from '../../context/AuthContext';
import { userService } from '../../Services/userService';
import { useNavigate } from 'react-router';

const AdminLogin = () => {
    const { login } = useAuth();
    const [formData, setFormData] = useState({
        username: '',
        password: '',
    });

    const navigate = useNavigate();
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const status = await userService.loginAdmin({
            mail: formData.username,
            password: formData.password
        })

        if (status === 200) {
            console.log("Admin Girişi Başarılı");
            login('Admin', true);
            navigate('/admin');
            return;
        }

        console.log("Bir şeyler ters gitti");
        // Şimdilik sahte admin girişi yapıyoruz
        //login('Admin', true);
    };

    return (
        <Container component="main" maxWidth="xs" sx={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Paper elevation={6} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', padding: 4, width: '100%', borderTop: '5px solid', borderColor: 'primary.main' }}>
                <Avatar sx={{ m: 1, bgcolor: 'primary.main', width: 56, height: 56 }}>
                    <LockOutlinedIcon fontSize="large" />
                </Avatar>
                <Typography component="h1" variant="h4" sx={{ fontWeight: 'bold', mb: 1, color: '#333' }}>
                    go-admin
                </Typography>
                <Typography component="h2" variant="subtitle1" sx={{ mb: 3, color: 'text.secondary' }}>
                    Yönetici Panel Girişi
                </Typography>
                <Box component="form" onSubmit={handleSubmit} noValidate sx={{ width: '100%' }}>
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        id="username"
                        label="Kullanıcı Adı"
                        name="username"
                        autoComplete="username"
                        autoFocus
                        value={formData.username}
                        onChange={handleChange}
                    />
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        name="password"
                        label="Şifre"
                        type="password"
                        id="password"
                        autoComplete="current-password"
                        value={formData.password}
                        onChange={handleChange}
                    />
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        size="large"
                        sx={{ mt: 3, mb: 2, py: 1.5, fontWeight: 'bold', textTransform: 'none', fontSize: '1.1rem' }}
                    >
                        Giriş Yap
                    </Button>
                </Box>
            </Paper>
        </Container>
    );
};

export default AdminLogin;
