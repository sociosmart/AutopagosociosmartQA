import * as React from "react";
import { styled, useTheme } from "@mui/material/styles";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Box from "@mui/material/Box";
import MuiDrawer from "@mui/material/Drawer";
import MuiAppBar from "@mui/material/AppBar";
import CampaignIcon from "@mui/icons-material/Campaign";
import Toolbar from "@mui/material/Toolbar";
import List from "@mui/material/List";
import CssBaseline from "@mui/material/CssBaseline";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import MenuIcon from "@mui/icons-material/Menu";
import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import DashboardIcon from "@mui/icons-material/Dashboard";
import SettingsIcon from "@mui/icons-material/Settings";
import EmojiEventsIcon from "@mui/icons-material/EmojiEvents";
import { matchPath, Outlet, useLocation, useNavigate } from "react-router-dom";

import { useGetMeQuery } from "../services/user";

import SmartGasLogo from "../assets/smartgas_logo.png";
import { CircularProgress } from "@mui/material";
import {
  AccountCircle,
  LocalGasStation,
  MapsHomeWork,
  Paid,
  CloudSync,
} from "@mui/icons-material";
import { logOut } from "../store/authSlice";
import { useDispatch } from "react-redux";

const drawerWidth = 240;

const openedMixin = (theme) => ({
  width: drawerWidth,
  transition: theme.transitions.create("width", {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen,
  }),
  overflowX: "hidden",
});

const closedMixin = (theme) => ({
  transition: theme.transitions.create("width", {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  overflowX: "hidden",
  width: `calc(${theme.spacing(7)} + 1px)`,
  [theme.breakpoints.up("sm")]: {
    width: `calc(${theme.spacing(8)} + 1px)`,
  },
});

const DrawerHeader = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "flex-end",
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
}));

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== "open",
})(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(["width", "margin"], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(["width", "margin"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Drawer = styled(MuiDrawer, {
  shouldForwardProp: (prop) => prop !== "open",
})(({ theme, open }) => ({
  width: drawerWidth,
  flexShrink: 0,
  whiteSpace: "nowrap",
  boxSizing: "border-box",
  ...(open && {
    ...openedMixin(theme),
    "& .MuiDrawer-paper": openedMixin(theme),
  }),
  ...(!open && {
    ...closedMixin(theme),
    "& .MuiDrawer-paper": closedMixin(theme),
  }),
}));

const routeToDisplay = {
  "/protected/users": "Usuarios",
  "/protected/gas-pumps": "Bombas",
  "/protected/gas-stations": "Gasolineras",
  "/protected/payments": "Pagos",
  "/protected": "Dashboard",
  "/protected/synchronizations": "Sincronizaciones",
  "/protected/synchronizations/:id/details": "Detalles de sincronizacion",
  "/protected/users/new": "Nuevo usuario",
  "/protected/users/:id/edit": "Editar usuario",
  "/protected/settings": "Configuraciones",
  "/protected/campaigns": "Campañas",
  "/protected/campaigns/:id/edit": "Editar campaña",
  "/protected/elegibility/levels": "Niveles",
  "/protected/elegibility/customers/levels": "Niveles de clientes",
};

const menu = [
  {
    icon: <DashboardIcon />,
    text: "Dashboard",
    navigateTo: "",
  },
  {
    icon: <AccountCircle />,
    text: "Usuarios",
    navigateTo: "/protected/users",
  },
  {
    icon: <MapsHomeWork />,
    text: "Estaciones",
    navigateTo: "/protected/gas-stations",
    requiredPermission: "view_gas_stations",
  },
  {
    icon: <LocalGasStation />,
    text: "Bombas",
    navigateTo: "/protected/gas-pumps",
    requiredPermission: "view_gas_pumps",
  },
  {
    icon: <Paid />,
    text: "Pagos",
    navigateTo: "/protected/payments",
    requiredPermission: "view_payments",
  },
  {
    icon: <CloudSync />,
    text: "Sincronizaciones",
    navigateTo: "/protected/synchronizations",
    requiredPermission: "view_synchronizations",
  },
  {
    icon: <SettingsIcon />,
    text: "Configuraciones",
    navigateTo: "/protected/settings",
  },
  {
    icon: <CampaignIcon />,
    text: "Campañas",
    navigateTo: "/protected/campaigns",
    requiredPermission: "view_campaigns",
  },
  {
    icon: <EmojiEventsIcon />,
    text: "Niveles",
    navigateTo: "/protected/elegibility/levels",
    requiredPermission: "view_elegibility_levels",
  },
  {
    icon: <EmojiEventsIcon />,
    text: "Niveles de clientes",
    navigateTo: "/protected/elegibility/customers/levels",
    requiredPermission: "view_customer_levels",
  },
];

export default function MiniDrawer() {
  const theme = useTheme();
  const [open, setOpen] = React.useState(true);
  const [anchorEl, setAnchorEl] = React.useState(null);

  const dispatch = useDispatch();

  const { data: user = {}, isLoading } = useGetMeQuery();

  const navigate = useNavigate();

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  const handleProfileMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const currentLocation = useLocation();

  const route = React.useMemo(
    () =>
      Object.keys(routeToDisplay).find((key) =>
        matchPath(key, currentLocation.pathname),
      ),
    [currentLocation.pathname],
  );

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const isMenuOpen = Boolean(anchorEl);
  const menuId = "primary-search-account-menu";
  const renderMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      id={menuId}
      keepMounted
      transformOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      {user?.is_admin && (
        <MenuItem
          onClick={() => {
            handleMenuClose();
            navigate(`/protected/users/${user?.id}/edit`);
          }}
        >
          Editar perfil
        </MenuItem>
      )}
      <MenuItem
        onClick={() => {
          dispatch(logOut());
        }}
      >
        LogOut
      </MenuItem>
    </Menu>
  );

  return (
    <Box sx={{ display: "flex", height: "100vh" }}>
      <CssBaseline />
      <AppBar position="absolute" open={open}>
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            onClick={handleDrawerOpen}
            edge="start"
            sx={{
              marginRight: 5,
              ...(open && { display: "none" }),
            }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div">
            {routeToDisplay[route]}
          </Typography>

          <Box sx={{ flexGrow: 1 }} />
          <Typography sx={{ display: { xs: "none", sm: "block" } }}>
            {user?.first_name} {user?.last_name}
          </Typography>
          <IconButton
            size="large"
            edge="end"
            aria-label="account of current user"
            aria-controls={menuId}
            aria-haspopup="true"
            onClick={handleProfileMenuOpen}
            color="inherit"
          >
            <AccountCircle />
          </IconButton>
        </Toolbar>
        {renderMenu}
      </AppBar>
      <Drawer variant="permanent" open={open}>
        <DrawerHeader>
          <Box maxWidth="100%">
            <img src={SmartGasLogo} width="100%" />
          </Box>
          <IconButton onClick={handleDrawerClose}>
            {theme.direction === "rtl" ? (
              <ChevronRightIcon />
            ) : (
              <ChevronLeftIcon />
            )}
          </IconButton>
        </DrawerHeader>
        <Divider />
        <List>
          {menu.map((m, idx) => {
            if (
              user?.is_admin ||
              user?.permissions?.find((p) => m.requiredPermission === p)
            )
              return (
                <ListItem
                  disablePadding
                  sx={{ display: "block" }}
                  key={idx}
                  onClick={() => {
                    navigate(m.navigateTo);
                  }}
                >
                  <ListItemButton
                    sx={{
                      minHeight: 48,
                      justifyContent: open ? "initial" : "center",
                      px: 2.5,
                    }}
                  >
                    <ListItemIcon
                      sx={{
                        minWidth: 0,
                        mr: open ? 3 : "auto",
                        justifyContent: "center",
                      }}
                    >
                      {m.icon}
                    </ListItemIcon>
                    <ListItemText
                      primary={m.text}
                      sx={{ opacity: open ? 1 : 0 }}
                    />
                  </ListItemButton>
                </ListItem>
              );
          })}
        </List>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <DrawerHeader />
        {isLoading ? (
          <Box
            minHeight="100%"
            sx={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              flexDirection: "column",
            }}
          >
            <CircularProgress />
            <Typography mt={3}>Cargando...</Typography>
          </Box>
        ) : (
          <Box
            sx={{ display: "flex", flexDirection: "column", height: "100%" }}
          >
            <Outlet />
          </Box>
        )}
      </Box>
    </Box>
  );
}
