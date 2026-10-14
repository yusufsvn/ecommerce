import React, { useState } from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Badge,
  Box,
  IconButton,
  InputBase,
  alpha,
  styled,
  Tooltip,
} from '@mui/material';
import ShoppingCartOutlinedIcon from '@mui/icons-material/ShoppingCartOutlined';
import PersonOutlineIcon from '@mui/icons-material/PersonOutline';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import SearchIcon from '@mui/icons-material/Search';
import { useNavigate } from 'react-router';
import type { NavbarProps } from '../models/Props';
import { useAuth } from '../context/AuthContext';

const Search = styled('div')(({ theme }) => ({
  position: 'relative',
  borderRadius: 9999, // Pill shape per DESIGN.md
  backgroundColor: '#F3F4F6',
  border: '1px solid #E5E7EB',
  transition: 'all 0.2s ease',
  '&:hover': {
    backgroundColor: '#EEF2F6',
    borderColor: '#D1D5DB',
  },
  '&:focus-within': {
    backgroundColor: '#FFFFFF',
    borderColor: '#2563EB',
    boxShadow: '0 0 0 3px rgba(37, 99, 235, 0.15)',
  },
  marginRight: theme.spacing(2),
  marginLeft: 0,
  width: '100%',
  maxWidth: '420px',
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: '100%',
  position: 'absolute',
  top: 0,
  pointerEvents: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: theme.palette.text.secondary,
  right: 0,
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: 'inherit',
  width: '100%',
  '& .MuiInputBase-input': {
    padding: '10px 16px',
    paddingRight: '42px',
    fontSize: '0.875rem',
    fontFamily: '"Inter", sans-serif',
    color: theme.palette.text.primary,
  },
}));

function Navbar(navbarProps: NavbarProps) {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearchKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && searchQuery.trim() !== '') {
      navigate(`/products?searchQuery=${encodeURIComponent(searchQuery.trim())}`);
      setSearchQuery('');
    }
  };

  return (
    <AppBar
      position="sticky"
      elevation={0}
      sx={{
        height: 72,
        bgcolor: '#FFFFFF',
        color: 'text.primary',
        borderBottom: '1px solid #E5E7EB',
        justifyContent: 'center',
        zIndex: (theme) => theme.zIndex.drawer + 1,
      }}
    >
      <Toolbar sx={{ height: 72, minHeight: '72px !important', px: { xs: 2, md: 4 }, justifyContent: 'space-between' }}>
        {/* Left Side: Logo & Links */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: { xs: 2, md: 3.5 } }}>
          <Typography
            variant="h6"
            onClick={() => navigate('/')}
            sx={{
              cursor: 'pointer',
              fontFamily: '"Outfit", sans-serif',
              fontWeight: 800,
              color: 'primary.main',
              letterSpacing: '-0.5px',
              fontSize: '1.25rem',
            }}
          >
            Ecommerce Store
          </Typography>

          <Box sx={{ display: { xs: 'none', md: 'flex' }, alignItems: 'center', gap: 1 }}>
            <Button
              color="inherit"
              onClick={() => navigate('/products')}
              sx={{
                fontWeight: 600,
                fontSize: '0.9rem',
                color: 'text.primary',
                '&:hover': { color: 'primary.main', bgcolor: 'transparent' },
              }}
            >
              Ürünler
            </Button>
            <Button
              color="inherit"
              onClick={() => navigate('/products')}
              sx={{
                fontWeight: 600,
                fontSize: '0.9rem',
                color: 'text.secondary',
                '&:hover': { color: 'primary.main', bgcolor: 'transparent' },
              }}
            >
              Kategoriler
            </Button>
            <Button
              color="inherit"
              onClick={() => navigate('/products')}
              sx={{
                fontWeight: 600,
                fontSize: '0.9rem',
                color: 'text.secondary',
                '&:hover': { color: 'primary.main', bgcolor: 'transparent' },
              }}
            >
              Kampanyalar
            </Button>
          </Box>
        </Box>

        {/* Center: Search Bar */}
        <Box sx={{ flexGrow: 1, display: 'flex', justifyContent: 'center', px: 2 }}>
          <Search>
            <StyledInputBase
              placeholder="Arama..."
              inputProps={{ 'aria-label': 'arama' }}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={handleSearchKeyPress}
            />
            <SearchIconWrapper>
              <SearchIcon fontSize="small" />
            </SearchIconWrapper>
          </Search>
        </Box>

        {/* Right Side: Icons */}
        <Box sx={{ display: 'flex', gap: 0.5, alignItems: 'center' }}>
          <Tooltip title="Profil">
            <IconButton
              color="inherit"
              onClick={() => (isAuthenticated ? navigate('/profile') : navigate('/login'))}
              sx={{ color: 'text.primary', '&:hover': { color: 'primary.main' } }}
            >
              <PersonOutlineIcon />
            </IconButton>
          </Tooltip>

          <Tooltip title="Favoriler">
            <IconButton
              color="inherit"
              onClick={() => navigate('/products')}
              sx={{ color: 'text.primary', '&:hover': { color: 'primary.main' } }}
            >
              <FavoriteBorderIcon />
            </IconButton>
          </Tooltip>

          <Tooltip title="Sepet">
            <IconButton
              color="inherit"
              onClick={() => navigate('/cart')}
              aria-label="sepet"
              sx={{ color: 'text.primary', '&:hover': { color: 'primary.main' } }}
            >
              <Badge
                badgeContent={navbarProps.getTotalItems()}
                sx={{
                  '& .MuiBadge-badge': {
                    bgcolor: 'secondary.main',
                    color: 'white',
                    fontWeight: 700,
                    fontSize: '0.75rem',
                  },
                }}
              >
                <ShoppingCartOutlinedIcon />
              </Badge>
            </IconButton>
          </Tooltip>
        </Box>
      </Toolbar>
    </AppBar>
  );
}

export default Navbar;

