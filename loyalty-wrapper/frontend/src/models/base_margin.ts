import { BaseProduct } from "../models/base_product";
import { BaseGasStation } from "../models/base_gas_station";

export interface BaseMargin {
    id: string;
    marginType: string;
    margin: number;
    points: number;
    product: BaseProduct;
    gasStation: BaseGasStation | null ;
}