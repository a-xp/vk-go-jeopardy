import {
    Avatar,
    Button,
    Caption,
    Cell,
    FormItem,
    FormLayout,
    Input,
    Panel,
    PanelHeader,
    PullToRefresh,
    View
} from "@vkontakte/vkui";
import {useCallback, useContext, useEffect, useRef, useState} from "react";
import api from "../service/api";
import {GlobalActions} from "../App";

function GroupsPage() {

    const [list, setList] = useState();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);
    const globalActions = useContext(GlobalActions);

    const onRefresh = useCallback(async () => {
        try {
            setLoading(true);
            const result = await api.getGroups();
            setList(result.items);
        } catch (e) {
            setError(true);
        }
        setLoading(false)
    }, [setList, setLoading, setError]);

    useEffect(() => {
        onRefresh();
    }, [onRefresh]);

    const onRemove = useCallback(async (id) => {
        try {
            await api.deleteGroup(id);
            return onRefresh();
        } catch (e) {
            globalActions.openSnack("Не удалось удалить сообщество")
        }
    }, [globalActions, onRefresh]);

    const apiKeyRef = useRef();

    const onAdd = useCallback(async () => {
        try {
            const payload = {
                apiKey: apiKeyRef.current.value,
            };
            await api.addGroup(payload);
            apiKeyRef.current.value = '';
            return onRefresh();
        } catch (e) {
            console.error(e);
            globalActions.openSnack("Не удалось добавить сообщество")
        }
    }, [globalActions, onRefresh, apiKeyRef]);

    return (
        <View id="settings" activePanel="settings">
            <Panel id="settings">
                <PanelHeader>Сообщества</PanelHeader>
                {error && <Caption weight="bold" level="1" className="centered-msg">Ошибка загрузки</Caption>}
                {!error && <PullToRefresh onRefresh={onRefresh} isFetching={loading}>
                    {!loading && list && list.length > 0 && list.map(e =>
                        <Cell key={e.id}
                              href={"https://vk.com/club" + e.id}
                              target="_blank"
                              rel="noopener noreferrer"
                              before={<Avatar src={e.image ? e.image : "https://vk.com/images/community_50.png"}
                                              size={48}/>}
                              removable
                              onRemove={() => onRemove(e.id)}>
                            {e.name}
                        </Cell>)}
                </PullToRefresh>}
                <FormLayout>
                    <FormItem top="Ключ доступа API">
                        <Input getRef={apiKeyRef}/>
                    </FormItem>
                    <FormItem>
                        <Button size="m" stretched onClick={onAdd}>Добавить сообщество</Button>
                    </FormItem>
                </FormLayout>
            </Panel>
        </View>
    )
}


export default GroupsPage;
