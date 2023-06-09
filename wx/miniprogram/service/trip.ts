
import { rental } from "./proto_gen/rental/rental_pb";
import { Rentcar } from "./request";

export namespace TripService{
    export function CreateTrip(req: rental.v1.ICreateTripRequest):Promise<rental.v1.ITripEntity>{
          return Rentcar.sendRequestWithAuthRetry({
            method:'POST',
            path:'/v1/trip',
            data:req,
            respMarshaller:rental.v1.TripEntity.fromObject,
        })
    }
}