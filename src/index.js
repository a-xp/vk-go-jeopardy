import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import "@vkontakte/vkui/dist/vkui.css";
import {AdaptivityProvider, ConfigProvider} from "@vkontakte/vkui";
import api from "./service/api";
import bridge from '@vkontakte/vk-bridge';
import 'core-js/es/map';
import 'core-js/es/set';

window.addEventListener('error', (message) => {
    api.sendError(message);
});

bridge.send("VKWebAppInit", {});

ReactDOM.render(
    <React.StrictMode>
        <ConfigProvider>
            <AdaptivityProvider>
                <App/>
            </AdaptivityProvider>
        </ConfigProvider>
    </React.StrictMode>,
    document.getElementById('root')
);

