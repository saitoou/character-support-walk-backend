import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    stages: [
        { duration: "30s", target: 30 },
        { duration: "1m", target: 100 },
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
    const params = {
        headers: {
            Authorization: `Bearer ${data.token}`,
            "Content-Type": "application/json",
        },
    };

    const homeRes = http.get(`${BASE_URL}/home`, params)
    check(homeRes, {
        "home status is 200": (r) => r.status === 200,
    });

    const walkOptionsRes = http.get(`${BASE_URL}/walk-options`, params)
    check(walkOptionsRes, {
        "walk options status is 200": (r) => r.status === 200,
    });

    const walkRes = http.post(`${BASE_URL}/walks`, JSON.stringify({
        walk_option_id: 1,
      }),
      params
    );
    check(walkRes, {
        "walks status is 200": (r) => r.status === 200,
        "has walk id": (r) => !!r.json("walk_id"),
    });

    const walkID = walkRes.json("walk_id");
    const completeRes = http.patch(`${BASE_URL}/walks/${walkID}/complete`, null, params);
    check(completeRes, {
        "complete walk status is 200": (r) => r.status === 200,
    });

    sleep(1);
}