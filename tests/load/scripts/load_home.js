import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    stages: [
        { duration: "30s", target: 10 },
        { duration: "1m", target: 30 },
        { duration: "30s", target: 50 },
        { duration: "30s", target: 0 },
    ],
    thresholds: {
        http_req_failed: ["rate<0.01"],
        http_req_duration: ["p(95)<800"],
    },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080/api/v1";

export function setup() {
    const res = http.post(`${BASE_URL}/auth/dev-login`);
    return {
        token: res.json("access_token"),
    };
}

export default function (data) {
    const res = http.get(`${BASE_URL}/home`, {
        headers: {
            Authorization: `Bearer ${data.token}`,
        },
    });

    check(res, {
        "home status is 200": (r) => r.status === 200,
    });
    
    if (res.status !== 200) {
        console.log(`status=${res.status}`);
        console.log(res.body);
      }

    sleep(0.2);
}