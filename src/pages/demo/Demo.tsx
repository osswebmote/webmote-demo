import useDeviceOrientation from "./useDeviceOrientation.tsx";
import {useEffect, useState} from "react";
import styles from "./styles.module.css";

const SVR_ADDR = "192.168.101.191:8000";

export default function Demo() {
    const {orientation, requestPermission, isPermissionGranted, isSupported} =
        useDeviceOrientation();
    const [isStart, setIsStart] = useState(false);
    const [ws, setWs] = useState<WebSocket | null>(null);

    useEffect(() => {
        if (isStart) {
            const lws = ws!;
            if (lws.readyState === WebSocket.OPEN) {
                lws.send(
                    JSON.stringify({
                        event: "data",
                        data: orientation,
                    })
                );
            } else {
                console.log("error");
            }
        }
    }, [orientation]);

    async function handleStart() {
        if (isStart) {
            console.log("Stop");
            ws!.close(1000);
            setIsStart(false);
            return;
        }
        console.log("Start");

        if (isSupported) {
            if (!isPermissionGranted) {
                await requestPermission();
            }
        }

        setIsStart(true);
        const localWs = new WebSocket(`wss://${SVR_ADDR}/ws/test`);
        setWs(localWs);
    }

    function handleCalibrate(isLeft: boolean, isTop: boolean) {
        if (!isStart) {
            return;
        }

        let event = "";

        if (isLeft) {
            if (isTop) {
                event = "lt";
            } else {
                event = "lb";
            }
        } else {
            if (isTop) {
                event = "rt";
            } else {
                event = "rb";
            }
        }

        ws!.send(
            JSON.stringify({
                event: event,
                data: orientation,
            })
        );
    }

    return (
        <div className={styles.center_div}>
            <p>{isStart ? "시작함!" : "멈춤!"}</p>
            <button className={styles.toggle} onClick={handleStart}>
                {isStart ? "멈추기" : "시작하기"}
            </button>
            <div>
                <p>alpha : {isStart ? orientation.alpha : 0}</p>
                <p>beta : {isStart ? orientation.beta : 0}</p>
                <p>gamma : {isStart ? orientation.gamma : 0}</p>
            </div>

            <button
                className={styles.calib}
                onClick={() => {
                    handleCalibrate(true, true);
                }}
            >
                Left Top Calibration
            </button>
            <button
                className={styles.calib}
                onClick={() => {
                    handleCalibrate(true, false);
                }}
            >
                Left Bottom Calibration
            </button>
            <button
                className={styles.calib}
                onClick={() => {
                    handleCalibrate(false, true);
                }}
            >
                Right Top Calibration
            </button>
            <button
                className={styles.calib}
                onClick={() => {
                    handleCalibrate(false, false);
                }}
            >
                Right Bottom Calibration
            </button>
        </div>
    );
}
