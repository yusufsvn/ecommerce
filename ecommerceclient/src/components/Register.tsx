import React, { useState } from 'react';
import { Container, TextField, Button, Typography, Box, Link, Paper, Grid, Stepper, Step, StepLabel } from '@mui/material';
import { useNavigate } from 'react-router';
import { userService } from '../Services/userService';
import type { Registeruser } from '../models/userModels';
import { useAuth } from '../context/AuthContext';

const steps = ['Kişisel Bilgiler', 'Güvenlik'];

const Register = () => {
    const navigate = useNavigate();
    const { login } = useAuth();

    const [activeStep, setActiveStep] = useState(0);

    const [formData, setFormData] = useState<Registeruser>({
        name: '',
        second_name: '',
        mail: '',
        password: '',
        password_check: ''
    });
    
    const [nameError, setNameError] = useState('');
    const [secondNameError, setSecondNameError] = useState('');
    const [emailError, setEmailError] = useState('');

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleNext = () => {
        let isValid = true;

        if (activeStep === 0) {
            if (!formData.name.trim()) {
                setNameError('Ad alanı zorunludur.');
                isValid = false;
            } else {
                setNameError('');
            }

            if (!formData.second_name.trim()) {
                setSecondNameError('Soyad alanı zorunludur.');
                isValid = false;
            } else {
                setSecondNameError('');
            }

            const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
            if (!emailRegex.test(formData.mail)) {
                setEmailError('Lütfen geçerli bir e-posta adresi giriniz.');
                isValid = false;
            } else {
                setEmailError('');
            }
        }

        // We could add more validations for other steps if necessary

        if (isValid) {
            setActiveStep((prevActiveStep) => prevActiveStep + 1);
        }
    };

    const handleBack = () => {
        setActiveStep((prevActiveStep) => prevActiveStep - 1);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        
        if (formData.password !== formData.password_check) {
            return;
        }

        const statusCode = await userService.registerUser(formData)
        if (statusCode == 201) {
            login(`${formData.name} ${formData.second_name}`);
            navigate('/')
        }
    };

    const renderStepContent = (step: number) => {
        switch (step) {
            case 0:
                return (
                    <Grid container spacing={2} justifyContent="center">
                        <Grid size={{ xs: 12, sm: 6 }}>
                            <TextField
                                required
                                fullWidth
                                id="name"
                                label="Ad"
                                name="name"
                                autoComplete="name"
                                autoFocus
                                value={formData.name}
                                onChange={handleChange}
                                error={!!nameError}
                                helperText={nameError}
                            />
                        </Grid>
                        <Grid size={{ xs: 12, sm: 6 }}>
                            <TextField
                                required
                                fullWidth
                                id="second_name"
                                label="Soyad"
                                name="second_name"
                                autoComplete="family-name"
                                value={formData.second_name}
                                onChange={handleChange}
                                error={!!secondNameError}
                                helperText={secondNameError}
                            />
                        </Grid>
                        <Grid size={{ xs: 12, sm: 6 }}>
                            <TextField
                                required
                                fullWidth
                                id="mail"
                                label="Email Adresi"
                                name="mail"
                                autoComplete="email"
                                value={formData.mail}
                                onChange={handleChange}
                                error={!!emailError}
                                helperText={emailError}
                            />
                        </Grid>
                    </Grid>
                );
            case 1:
                return (
                    <Grid container spacing={2} justifyContent="center">
                        <Grid size={{ xs: 12, sm: 6 }}>
                            <TextField
                                required
                                fullWidth
                                name="password"
                                label="Şifre"
                                type="password"
                                id="password"
                                autoComplete="new-password"
                                value={formData.password}
                                onChange={handleChange}
                            />
                        </Grid>
                        <Grid size={{ xs: 12, sm: 6 }}>
                            <TextField
                                required
                                fullWidth
                                name="password_check"
                                label="Şifre Tekrar"
                                type="password"
                                id="password_check"
                                autoComplete="new-password"
                                value={formData.password_check}
                                onChange={handleChange}
                                error={formData.password !== '' && formData.password_check !== '' && formData.password !== formData.password_check}
                                helperText={formData.password !== '' && formData.password_check !== '' && formData.password !== formData.password_check ? 'Şifreler uyuşmuyor' : ''}
                            />
                        </Grid>
                    </Grid>
                );
            default:
                return 'Bilinmeyen adım';
        }
    };

    return (
        <Container component="main" maxWidth="md">
            <Paper elevation={6} sx={{ marginTop: 8, marginBottom: 8, display: 'flex', flexDirection: 'column', alignItems: 'center', padding: { xs: 4, md: 6 }, borderRadius: 4 }}>
                <Typography component="h1" variant="h4" sx={{ fontWeight: 'bold', mb: 1, color: 'primary.main' }}>
                    Kayıt Ol
                </Typography>
                <Typography variant="body1" color="text.secondary" sx={{ mb: 4, textAlign: 'center' }}>
                    Yeni bir hesap oluşturmak için bilgilerinizi eksiksiz giriniz.
                </Typography>
                
                <Stepper activeStep={activeStep} sx={{ width: '100%', mb: 4 }} alternativeLabel>
                    {steps.map((label) => (
                        <Step key={label}>
                            <StepLabel>{label}</StepLabel>
                        </Step>
                    ))}
                </Stepper>

                <Box component="form" onSubmit={handleSubmit} noValidate sx={{ width: '100%' }}>
                    {renderStepContent(activeStep)}
                    
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 4, mb: 2 }}>
                        <Button
                            color="inherit"
                            disabled={activeStep === 0}
                            onClick={handleBack}
                            sx={{ mr: 1, borderRadius: 2 }}
                            variant="outlined"
                        >
                            Geri
                        </Button>
                        {activeStep === steps.length - 1 ? (
                            <Button
                                type="submit"
                                variant="contained"
                                sx={{ borderRadius: 2 }}
                                disabled={formData.password !== formData.password_check || !formData.password}
                            >
                                Kayıt Ol
                            </Button>
                        ) : (
                            <Button
                                variant="contained"
                                onClick={handleNext}
                                sx={{ borderRadius: 2 }}
                            >
                                İleri
                            </Button>
                        )}
                    </Box>
                    <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
                        <Link href="#" variant="body2" onClick={(e) => { e.preventDefault(); navigate('/login'); }} sx={{ textDecoration: 'none', fontWeight: 600, '&:hover': { textDecoration: 'underline' } }}>
                            {"Zaten hesabın var mı? Giriş Yap"}
                        </Link>
                    </Box>
                </Box>
            </Paper>
        </Container>
    );
};

export default Register;
