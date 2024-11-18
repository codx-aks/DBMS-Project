import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

export const options = {
    vus: 1, 
    duration: '3m', 
};

export default function () {
    const startRollNo = 107100001;
    const endRollNo = 1071100000;

    for (let rollNo = startRollNo; rollNo <= endRollNo; rollNo++) {
        const data = JSON.stringify({
            roll_no: rollNo.toString(),
            vendor_id: "1021603896624021505", 
            amount: 25,
            pin: "1234", 
        });

        const headers = {
            Cookie: "session_token=jGNVfAbSOi0b6xVbAkboq2qxqodUrsAnlWOpDN8M7Xo",
        };

        const paramsWrite = {
            headers: {
                Cookie: headers.Cookie,
                "Content-Type": "application/json",
            },
        };

        let res = http.post('http://localhost:7070/pay', data, paramsWrite);
        check(res, { 'success payment': (r) => r.status === 200 });


    }
}
