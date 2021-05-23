import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import "@vkontakte/vkui/dist/vkui.css";
import {AdaptivityProvider, ConfigProvider} from "@vkontakte/vkui";

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

