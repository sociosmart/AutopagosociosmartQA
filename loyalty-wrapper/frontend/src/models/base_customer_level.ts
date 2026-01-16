import { BaseCustomer } from "../models/base_customer";
import { BaseLevel } from "../models/base_level";

export interface BaseCustomerLevel {
    id: string;
    customer: BaseCustomer;
    level: BaseLevel;
    startDate: string;
    endDate: string;
}