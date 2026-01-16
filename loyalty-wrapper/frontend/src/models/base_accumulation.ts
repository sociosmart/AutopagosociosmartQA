import { BaseGasStation } from './base_gas_station';
import { BaseCustomer } from './base_customer'; 
import { BaseProduct } from './base_product';

export interface BaseAccumulation {
    id: string;
    margin: number;
    amount: number;
    points: number;
    generatedPoints: number;
    marginType: 'by_margin' | 'by_liter';
    usedPoints: number;
    createdAt: string;
    gasStation: BaseGasStation;
    customer: BaseCustomer;
    product: BaseProduct;
    gasPrice: number;
}