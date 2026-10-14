import { Outlet, useLocation, useNavigate } from "react-router"
import { AppBar, Box, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography, Button } from "@mui/material"
import DashboardIcon from '@mui/icons-material/Dashboard';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import InventoryIcon from '@mui/icons-material/Inventory';
import LocalShippingIcon from '@mui/icons-material/LocalShipping';
import CategoryIcon from '@mui/icons-material/Category';
import AdminLogin from './AdminLogin';
import { useAuth } from '../../context/AuthContext';

const drawerWidth = 240;

const menuItems = [
  { text: 'Panel', icon: <DashboardIcon />, path: '/admin' },
  { text: 'Ürünler', icon: <InventoryIcon />, path: '/admin/products' },
  { text: 'Kategoriler', icon: <CategoryIcon />, path: '/admin/categories' },
  { text: 'Ürün Ekle', icon: <AddCircleIcon />, path: '/admin/create-product' },
  { text: 'Siparişler', icon: <LocalShippingIcon />, path: '/admin/orders' },
];

function Admin() {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAdmin, logout } = useAuth();

  if (!isAdmin) {
    return <AdminLogin />;
  }

  return (
    <Box sx={{ display: 'flex' }}>
      {/* Top App Bar */}
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1, borderTopLeftRadius: 0, borderTopRightRadius: 0 }}>
        <Toolbar sx={{ justifyContent: 'space-between' }}>
          <Typography variant="h6" noWrap component="div" sx={{ color: 'white', fontWeight: 700 }}>
            Admin Console
          </Typography>
          <Button color="inherit" onClick={() => { logout(); navigate('/admin') }}>
            Çıkış Yap
          </Button>
        </Toolbar>
      </AppBar>

      {/* Side Navigation (Drawer) */}
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' },
        }}
      >
        <Toolbar /> {/* Spacer for Top Bar */}
        <Box sx={{ overflow: 'auto' }}>
          <List>
            {menuItems.map((item) => {
              const isSelected = location.pathname === item.path;
              return (
                <ListItem key={item.text} disablePadding sx={{ px: 2, py: 0.5 }}>
                  <ListItemButton
                    onClick={() => navigate(item.path)}
                    selected={isSelected}
                    sx={{
                      borderRadius: 2, // Smooth edges for list items
                      '&.Mui-selected': {
                        bgcolor: 'primary.light',
                        color: 'primary.contrastText',
                        '& .MuiListItemIcon-root': {
                          color: 'primary.contrastText',
                        },
                        '&:hover': {
                          bgcolor: 'primary.main',
                        }
                      }
                    }}
                  >
                    <ListItemIcon sx={{ color: isSelected ? 'inherit' : 'text.secondary' }}>
                      {item.icon}
                    </ListItemIcon>
                    <ListItemText primary={item.text} primaryTypographyProps={{ fontWeight: isSelected ? 600 : 400 }} />
                  </ListItemButton>
                </ListItem>
              );
            })}
          </List>
        </Box>
      </Drawer>

      {/* Main Content Area */}
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Toolbar /> {/* Spacer so content isn't hidden behind AppBar */}
        <Outlet />
      </Box>
    </Box>
  )
}

export default Admin