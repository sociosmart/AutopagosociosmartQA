import { BaseCustomer } from "./base_customer";

export interface BaseReport {
    id : string;
    totalTransactions : number;
    avgAmount : number;
    totalAmount : number;
    totalGeneratedPoints : number;
    totalUsedPoints : number;
    totalPoints : number;
    customer : BaseCustomer;
  }


