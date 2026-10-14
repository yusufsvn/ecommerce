import { createTheme, CssBaseline } from '@mui/material';
import './App.css'
import { ThemeProvider } from '@emotion/react'
import { Route, Routes } from 'react-router';
import ProductDetail from './components/ProductDetail';
import Notfoundpage from './components/Notfoundpage';
import Home from './components/Home';
import ProductsPage from './components/ProductsPage';
import About from './components/About';
import Contact from './components/Contact';
import { useEffect, useState, useRef } from 'react';
import type { cartItem } from './models/cartmodel';
import Cart from './components/Cart';
import Admin from './components/admin/Admin';
import PublicLayout from './components/PublicLayout';
import CreateProduct from './components/admin/CreateProduct';
import AdminProducts from './components/admin/AdminProducts';
import EditProduct from './components/admin/Pop-ups/Edit';
import OrderSummary from './components/OrderSummary';
import Login from './components/Login';
import Register from './components/Register';
import Profile from './components/Profile';
import UserOrders from './components/user/UserOrders';
import UserAddress from './components/user/UserAddress';
import AdminDashboard from './components/admin/AdminDashboard';
import AdminOrders from './components/admin/AdminOrders';
import AdminCategories from './components/admin/AdminCategories';
import { cartService } from './Services/cartService';
import { useAuth } from './context/AuthContext';
import type { UpdateCartItem } from './models/ProductResponse';
import { toast } from 'react-toastify';


const theme = createTheme({
  palette: {
    primary: {
      main: '#1E3A8A', // Structural baseline, headers, primary actions
      light: '#2563EB',
      dark: '#172554',
      contrastText: '#FFFFFF',
    },
    secondary: {
      main: '#F97316', // High-conversion CTA, limited-time offers
      light: '#FB923C',
      dark: '#EA580C',
      contrastText: '#FFFFFF',
    },
    info: {
      main: '#38BDF8', // Metadata chips, tags, informational alerts
      contrastText: '#121C2A',
    },
    background: {
      default: '#F3F4F6', // Flat canvas layer (Level 0)
      paper: '#FFFFFF',   // Card & surface layer (Level 1)
    },
    text: {
      primary: '#374151',   // High readability dark gray
      secondary: '#6B7280', // Secondary gray
      disabled: '#9CA3AF',
    },
  },
  shape: {
    borderRadius: 10, // 10px default
  },
  typography: {
    fontFamily: '"Inter", sans-serif',
    h1: { fontFamily: '"Outfit", sans-serif', fontSize: '2.5rem', fontWeight: 700, color: '#1E3A8A' },
    h2: { fontFamily: '"Outfit", sans-serif', fontSize: '2rem', fontWeight: 700, color: '#1E3A8A' },
    h3: { fontFamily: '"Outfit", sans-serif', fontSize: '1.75rem', fontWeight: 600, color: '#1E3A8A' },
    h4: { fontFamily: '"Outfit", sans-serif', fontSize: '1.5rem', fontWeight: 600, color: '#1E3A8A' },
    h5: { fontFamily: '"Outfit", sans-serif', fontSize: '1.25rem', fontWeight: 600, color: '#1E3A8A' },
    h6: { fontFamily: '"Outfit", sans-serif', fontSize: '1.125rem', fontWeight: 600, color: '#1E3A8A' },
    body1: { fontFamily: '"Inter", sans-serif', fontSize: '1rem', fontWeight: 400, color: '#374151' },
    body2: { fontFamily: '"Inter", sans-serif', fontSize: '0.875rem', fontWeight: 400, color: '#6B7280' },
    button: {
      fontFamily: '"Inter", sans-serif',
      textTransform: 'none',
      fontWeight: 600,
    }
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 10, // DESIGN.md 4.2 (10px)
          boxShadow: 'none',
          textTransform: 'none',
          fontWeight: 600,
          '&:hover': {
            boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
          }
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12, // DESIGN.md 3.1 & shape.borderRadius.md: 12px
          boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.05), 0 1px 2px -1px rgba(0, 0, 0, 0.05)',
          border: '1px solid rgba(229, 231, 235, 0.8)',
        }
      }
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        }
      }
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: 10,
        }
      }
    }
  }
});

function App() {
  const { isAuthenticated, isAdmin } = useAuth();

  const [CartProducts, SetCartProducts] = useState<cartItem[]>([]);
  // Debounce için timer referanslarını tutacağımız obje
  const updateTimeoutRef = useRef<{ [key: string]: ReturnType<typeof setTimeout> }>({});

  //load cart from api
  useEffect(() => {
    if (isAuthenticated && !isAdmin) {
      cartService.getandProcessCartItems().then((cartItems: cartItem[]) => {
        SetCartProducts(cartItems);
      });
    } else {
      // Uyarıyı (cascading renders) çözmek için asenkron (gecikmeli) olarak state'i sıfırlıyoruz.
      Promise.resolve().then(() => SetCartProducts([]));
    }
  }, [isAdmin, isAuthenticated]);
  // cart functions
  // ADD TO CART
  const addtocart = (newProduct: cartItem) => {
    const qtyToAdd = newProduct.quantity && newProduct.quantity > 0 ? newProduct.quantity : 1;
    const TempCartProducts: cartItem[] = [...CartProducts];
    // Control of Have the product
    const TempProduct = TempCartProducts.find(item => item.product_id == newProduct.product_id)

    // Ürün Sepette Bulunmakta
    if (TempProduct) {
      // temp listesinde değişiklik yapılıyor
      const updatedItems = TempCartProducts.map((cart => {
        if (cart.product_id == newProduct.product_id) {
          const newQty = cart.quantity + qtyToAdd;
          const updatedata: UpdateCartItem = {
            product_id: newProduct.product_id,
            quantity: newQty
          }
          // api update product amount
          cartService.updateCartItemAmount(updatedata);
          return { ...cart, quantity: newQty }

        }
        return cart
      }))
      // set to SetCart
      SetCartProducts([...updatedItems])
      toast.success("Ürün miktarı güncellendi!");
      // api add to cart new product
    }
    // Ürün Sepette Değil
    else {
      console.log(`product ide debug: ${newProduct.product_id}`)
      cartService.addCartItem(newProduct.product_id, qtyToAdd);

      SetCartProducts([...TempCartProducts, { ...newProduct, quantity: qtyToAdd }])
      toast.success("Ürün sepete eklendi!");
    }
  }
  const deletefromcart = (id: string) => {
    console.log("delete func gelen id" + id)
    const TempCartProducts: cartItem[] = [...CartProducts];
    const newCartList = TempCartProducts.filter((cart) => cart.product_id !== id)
    SetCartProducts([...newCartList])
    cartService.deleteCartItem(id);
    toast.info("Ürün sepetten silindi.");
  }
  const getTotalItems = () => {

    let count: number = 0;

    for (let i = 0; i < CartProducts.length; i++) {
      count += CartProducts[i].quantity
    }

    return count;
  }

  const calculateTotalPrice = (): number => {
    const TempProducts: cartItem[] = [...CartProducts]

    let totalPrice: number = 0;
    TempProducts.forEach((element) => {
      totalPrice += (element.quantity * element.price)
    })
    return totalPrice;
  }
  const handleDecrement = (id: string) => {

    const TempCartProducts: cartItem[] = [...CartProducts];

    const handled: cartItem[] = TempCartProducts
      .map(item => {
        // Aranan ürünse miktarını azalt, değilse aynen bırak
        if (item.product_id === id) {
          const newQuantity = item.quantity - 1;

          // Önceki bekleyen bir API isteği varsa iptal et
          if (updateTimeoutRef.current[id]) {
            clearTimeout(updateTimeoutRef.current[id]);
          }

          if (newQuantity === 0) {
            // Miktar 0 olduysa ürünü sepetten tamamen sil (beklemeye gerek yok)
            cartService.deleteCartItem(id);
          } else {
            // 500ms bekleyerek son miktarı API'ye yolla (Debounce)
            updateTimeoutRef.current[id] = setTimeout(() => {
              cartService.updateCartItemAmount({ product_id: id, quantity: newQuantity });
            }, 750);
          }

          return { ...item, quantity: newQuantity };
        }
        return item;
      })
      .filter(item => item.quantity > 0);
    SetCartProducts([...handled])
  }

  const handleIncrement = (id: string) => {
    const TempCartProducts: cartItem[] = [...CartProducts];

    // .map() dizisini return ediyoruz ki state düzgün ezilsin (Eski hatalı kodu düzelttim)
    const handled = TempCartProducts.map((cartItem: cartItem) => {
      if (cartItem.product_id == id) {
        const newQuantity = cartItem.quantity + 1;

        // Önceki bekleyen bir API isteği varsa iptal et
        if (updateTimeoutRef.current[id]) {
          clearTimeout(updateTimeoutRef.current[id]);
        }

        // 500ms bekleyerek son miktarı API'ye yolla (Debounce)
        updateTimeoutRef.current[id] = setTimeout(() => {
          cartService.updateCartItemAmount({ product_id: id, quantity: newQuantity });
        }, 750);

        return { ...cartItem, quantity: newQuantity };
      }
      return cartItem;
    });

    SetCartProducts([...handled])
  }

  const clearCart = () => {
    SetCartProducts([]);
  }


  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {/*-----------Müşteri Paneli------------------*/}
      <Routes>
        {/*-----------Public Pages (Wrapped in PublicLayout)------------------*/}
        <Route element={<PublicLayout getTotalItems={getTotalItems} />}>
          <Route path='/' element={<Home addtocart={addtocart} cartItems={CartProducts} />} />
          <Route path='/products' element={<ProductsPage addtocart={addtocart} cartItems={CartProducts} />} />
          <Route path='/hakkimizda' element={<About />} />
          <Route path='/iletisim' element={<Contact />} />
          <Route path="/product/:id" element={<ProductDetail addtocart={addtocart} cartItems={CartProducts} />} />
          <Route path="/cart" element={<Cart products={CartProducts}
            handleIncrement={handleIncrement} handleDecrement={handleDecrement} calculateTotalPrice={calculateTotalPrice} deleteProduct={deletefromcart} />} />
          <Route path="/order-summary" element={<OrderSummary cartItems={CartProducts} totalPrice={calculateTotalPrice()} clearCart={clearCart} />} />
          
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/profile/orders" element={<UserOrders />} />
          <Route path="/profile/address" element={<UserAddress />} />
          <Route path="*" element={<Notfoundpage />} />
        </Route>

        {/*-----------Admin Panel------------------*/}
        <Route path='/admin' element={<Admin />}>
          <Route index element={<AdminDashboard />} />
          <Route path='products' element={<AdminProducts />} />
          <Route path='categories' element={<AdminCategories />} />
          <Route path='create-product' element={<CreateProduct />} />
          <Route path='edit-product/:id' element={<EditProduct />} />
          <Route path='orders' element={<AdminOrders />} />
        </Route>
      </Routes>
    </ThemeProvider>

  )
}

export default App
