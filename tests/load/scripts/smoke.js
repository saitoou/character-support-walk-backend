import http from "k6/http"
import { check, sleep } from "k6"

export const options = {
    vus: 1,
    iterations: 1,
    thresholds: {
        http_req_failed: ["rate<0.01"],
        http_req_duration: ["p(95)<500"],
    },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080/api/v1";

export default function () {
    const loginRes = http.post(`${BASE_URL}/auth/dev-login`);

    check(loginRes, {
        "dev login status is 200": (r) => r.status === 200,
        "has access token": (r) => !!r.json("access_token"),
    });

    const token = loginRes.json("access_token");

    const params = {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    };

    const homeRes = http.get(`${BASE_URL}/home`, params);

    check(homeRes, {
        "home status is 200": (r) => r.status === 200,
    });

    sleep(1);
}