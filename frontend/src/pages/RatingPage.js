import {Avatar, Caption, Cell, Panel, PanelHeader, PullToRefresh, View} from "@vkontakte/vkui";
import {useCallback, useEffect, useState} from "react";
import {formatRusNumerics} from "../service/util";
import api from "../service/api";

function RatingPage() {

    const [loading, setLoading] = useState(true);
    const [list, setList] = useState();
    const [userRating, setUserRating] = useState();
    const [error, setError] = useState();
    const [name, setName] = useState();

    const onRefresh = useCallback(async () => {
        try {
            setLoading(true);
            setError(false);
            const result = await api.getRating();
            setList(result.rating);
            setUserRating(result.userRating);
            setName(result.name);
        } catch (e) {
            console.error(e);
            setError(true);
        }
        setLoading(false);
    }, [setLoading, setError, setList, setUserRating]);

    useEffect(() => {
        onRefresh();
    }, []);

    if (error) {
        return <Caption weight="bold" level="1" className="centered-msg">Ошибка приложения</Caption>
    }

    return <View id="rating" activePanel="rating">
        <Panel id="rating">
            <PanelHeader>Рейтинг {name}</PanelHeader>
            <PullToRefresh onRefresh={onRefresh} isFetching={loading}>
                {!loading && list && (list.length > 0)
                && list.map(e => <Cell
                    key={e.pos}
                    href={"https://vk.com/id" + e.userId}
                    target="_blank"
                    rel="noopener noreferrer"
                    before={<Avatar src={e.img ? e.img : "https://vk.com/images/camera_50.png?ava=1"} size={48}/>}
                    description={e.score + " " + formatRusNumerics(e.score, ["балл", "балла", "баллов"])}>
                    <b>{e.pos}.</b> {e.name} {e.lastname}
                </Cell>)}
                {!loading && (!list || list.length === 0) &&
                <Caption weight="bold" level="1" className="centered-msg">Пока ничего нет</Caption>}
            </PullToRefresh>

        </Panel>
    </View>
}

export default RatingPage;
