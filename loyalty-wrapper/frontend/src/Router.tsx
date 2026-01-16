import { Routes, Route } from "react-router";
import Layout from "./components/Layout";
import HomeScreen from "./screens/Home";
import LoginScreen from "./screens/Login";
import CustomersScreen from "./screens/Customers";
import UsersScreen from "./screens/Users";
import GasStationsScreen from "./screens/GasStations";
import ProductsScreen from "./screens/Products";
import ProtectedRoute from "./components/ProtectedRoute"; 
import { Navigate } from "react-router-dom";
import  ProductsEditCreateScreen from "./screens/ProductsEditCreate";
import MarginsScreen from "./screens/Margins";
import MarginsEditCreateScreen from "./screens/MarginsEditCreate";
import AcumulationsScreen from "./screens/Accumulations";
import ReportScreen from "./screens/Report";
import LevelsScreen from "./screens/Levels";
import LevelsEditCreateScreen from "./screens/LevelsEditCreate";
import CustomerLevelsScreen from "./screens/CustomerLevels"
import BenefitsScreen from "./screens/Benefits"
import BenefitsEditCreateScreen from "./screens/BenefitsEditCreate"
import BenefitsGeneratedEditCreateScreen from "./screens/BenefitsGeneratedEditCreate"
import BenefitsGeneratedScreen from "./screens/BenefitsGenerated"
import UsersEditCreateScreen from "./screens/UsersEditCreate";
import BenefitsTicketsScreen from "./screens/BenefitsTicekts"

export default function Router() {
  return (
    <Routes>
      <Route path="login" element={<LoginScreen />} />


      <Route element={<ProtectedRoute />}>
        <Route path="protected" element={<Layout />}>
          <Route path="home" element={<HomeScreen />} />
          <Route path="customers" element={<CustomersScreen />} />
          <Route path="users" element={<UsersScreen />} />
          <Route path="gas-stations" element={<GasStationsScreen />} />
          <Route path="products" element={<ProductsScreen />} />
          <Route path="acumulations" element={<AcumulationsScreen />} />
          <Route path="margins" element={<MarginsScreen />} />
          <Route path="report" element={<ReportScreen />} />
          <Route path="levels" element={<LevelsScreen />} />
          <Route path="customer-levels" element={<CustomerLevelsScreen />} />
          <Route path="benefits" element={<BenefitsScreen />} />
          <Route path="benefits-generated" element={<BenefitsGeneratedScreen />} />
          <Route path="benefits-tickets" element={<BenefitsTicketsScreen />} />

          <Route path="products/create" element={<ProductsEditCreateScreen />} />   
          <Route path="margins/create" element={<MarginsEditCreateScreen />} />   
          <Route path="levels/create" element={<LevelsEditCreateScreen />} />
          <Route path="benefits/create" element={<BenefitsEditCreateScreen />} />
          <Route path="benefits-generated/create" element={<BenefitsGeneratedEditCreateScreen />} />
          <Route path="users/create" element={<UsersEditCreateScreen />}  />   
  
        </Route>

      </Route>
      <Route path="*" element={<Navigate to="/protected/users" />} />
      
    </Routes>
  );
}
