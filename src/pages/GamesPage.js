import {View} from "@vkontakte/vkui";
import {useCallback, useContext, useEffect, useState} from "react";
import {GlobalActions} from "../App";
import GamesList from "./GamesList";
import api from "../service/api";
import GameEditWrapper from "./GameEditWrapper";
import {DEFAULT_MESSAGES} from "../service/default";

function GamesPage() {

    const [currentGame, setCurrentGame] = useState(null);
    const [list, setList] = useState();
    const [loadingList, setLoadingList] = useState(true);
    const [editLvl, setEditLvl] = useState(0);
    const globalActions = useContext(GlobalActions);

    const onReload = useCallback(async () => {
        setLoadingList(true);
        try {
            const result = await api.getGames();
            setList(result.items);
        } catch (e) {
            globalActions.openSnack('Ошибка загрузки');
        }
        setLoadingList(false);
    }, [setList, setLoadingList, globalActions]);

    const onOpen = useCallback(async (id, cmd) => {
        setLoadingList(true);
        try {
            const game = await api.getGame(id);
            setCurrentGame(game);
            setEditLvl(cmd);
        } catch (e) {
            globalActions.openSnack('Ошибка загрузки');
        }
        setLoadingList(false);
    }, [setCurrentGame, setEditLvl, setLoadingList, globalActions]);

    const onClose = useCallback(() => {
        setCurrentGame(null);
    }, [setCurrentGame]);

    const onSave = useCallback(async (game) => {
        try {
            await api.updateGame(game);
            setCurrentGame(null);
            await onReload();
        } catch (e) {
            globalActions.openSnack('Ошибка сохранения');
        }
    }, [globalActions, setCurrentGame]);

    const onNew = useCallback(() => {
        setCurrentGame({
            name: 'Новая викторина',
            messages: DEFAULT_MESSAGES,
            active: false,
            new: true,
            post: {},
            rules: {},
            topics: []
        });
    }, [setCurrentGame]);

    useEffect(async () => {
        await onReload();
    }, [onReload]);

    return (
        <View activePanel="games">
            {!currentGame &&
            <GamesList id="games" list={list} loading={loadingList} onRefresh={onReload} onOpen={onOpen}
                       onNew={onNew}/>}
            {currentGame &&
            <GameEditWrapper id="games" game={currentGame} onClose={onClose} onSave={onSave} mode={editLvl}/>}
        </View>
    )
}

export default GamesPage;