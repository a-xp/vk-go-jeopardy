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

function AdminPage() {

    const [list, setList] = useState();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);

    const globalActions = useContext(GlobalActions);

    const onRefresh = useCallback(async () => {
        setLoading(true);
        setError(false);
        try {
            const result = await api.getAdmins();
            setList(result.items);
        } catch (e) {
            console.error(e);
            setError(true);
        }
        setLoading(false);
    }, [setLoading, setError, setList]);

    const onRemove = useCallback(async (id) => {
        try {
            await api.deleteAdmin(id);
            return onRefresh();
        } catch (e) {
            console.error(e);
            globalActions.openSnack('Ошибка при удалении');
        }
    }, []);

    const linkRef = useRef();

    const onAdd = useCallback(async () => {
        const value = linkRef.current.value;
        try {
            await api.addAdmin(value);
            linkRef.current.value = '';
            return onRefresh();
        } catch (e) {
            console.error(e);
            globalActions.openSnack('Не удалось добавить');
        }
    }, []);

    useEffect(() => {
        onRefresh();
    }, []);

    return (
        <View id="settings" activePanel="settings">
            <Panel id="settings">
                <PanelHeader>
                    Админы</PanelHeader>
                {!error && <PullToRefresh onRefresh={onRefresh} isFetching={loading}>
                    {!loading && list && list.length > 0 && list.map(e =>
                        <Cell key={e.id}
                              href={"https://vk.com/id" + e.userId}
                              target="_blank"
                              rel="noopener noreferrer"
                              removable
                              onRemove={() => onRemove(e.id)}
                              before={<Avatar src={e.image ? e.image : "https://vk.com/images/camera_50.png?ava=1"}
                                              size={48}/>}
                              data-id={e.id}>
                            {e.name}
                        </Cell>)}
                </PullToRefresh>}
                {error && <Caption weight="bold" level="1" className="centered-msg">Ошибка загрузки</Caption>}
                <FormLayout>
                    <FormItem top="Ссылка или ID">
                        <Input getRef={linkRef}/>
                    </FormItem>
                    <FormItem>
                        <Button size="m" stretched onClick={onAdd}>Добавить администратора</Button>
                    </FormItem>
                </FormLayout>
            </Panel>
        </View>
    )

}

export default AdminPage;