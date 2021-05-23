import {AppRoot, Epic, ScreenSpinner, Snackbar, Tabbar, TabbarItem, Title} from "@vkontakte/vkui";
import {createContext, useEffect, useMemo, useState} from "react";
import {
    Icon12Cancel,
    Icon28BrainOutline,
    Icon28FavoriteOutline,
    Icon28Users3Outline,
    Icon28UserStarBadgeOutline
} from "@vkontakte/icons";
import RatingPage from "./pages/RatingPage";
import api from "./service/api";
import AdminPage from "./pages/AdminPage";
import GroupsPage from "./pages/GroupsPage";
import GamesPage from "./pages/GamesPage";

export const GlobalActions = createContext(null);

function App() {
    const [activeStory, setActiveStory] = useState('rating');
    const onStoryChange = (e) => setActiveStory(e.currentTarget.dataset.story);
    const [isAdmin, setIsAdmin] = useState(false);
    const [publicUrl, setPublicUrl] = useState();
    const [error, setError] = useState();
    const [loading, setLoading] = useState(true);
    const [snack, setSnack] = useState(null);

    useEffect(async () => {
        try {
            const result = await api.getProfile();
            setIsAdmin(result.isAdmin);
            setPublicUrl(result.publicUrl);
        } catch (e) {
            setError(true)
        }
        setLoading(false);
    }, []);

    const context = useMemo(() => {
        return {
            openSnack(text) {
                setSnack(<Snackbar
                    before={<Icon12Cancel fill="red" width="16" height="16"/>}
                    onClose={() => setSnack(null)}>{text}</Snackbar>)
            },
            copyGameUrl(id) {
                const url = publicUrl + '#' + id;
                return navigator.clipboard.writeText(url);
            }
        }
    }, [setSnack, publicUrl]);

    if (error) {
        return <Title level="2" weight="heavy">Этим приложением можно пользоваться только в VK Mini Apps</Title>;
    }

    return (
        <GlobalActions.Provider value={context}>
            <AppRoot>
                {loading && <ScreenSpinner/>}
                {!loading && !error && isAdmin && <Epic activeStory={activeStory} tabbar={<Tabbar>
                    <TabbarItem onClick={onStoryChange} selected={activeStory === 'rating'} data-story="rating"
                                text="Рейтинг">
                        <Icon28FavoriteOutline/>
                    </TabbarItem>
                    <TabbarItem onClick={onStoryChange} selected={activeStory === 'games'} data-story="games"
                                text="Викторины">
                        <Icon28BrainOutline/>
                    </TabbarItem>
                    <TabbarItem onClick={onStoryChange} selected={activeStory === 'groups'} data-story="groups"
                                text="Сообщества">
                        <Icon28Users3Outline/>
                    </TabbarItem>
                    <TabbarItem onClick={onStoryChange} selected={activeStory === 'admins'} data-story="admins"
                                text="Админы">
                        <Icon28UserStarBadgeOutline/>
                    </TabbarItem>
                </Tabbar>}>
                    <RatingPage id="rating"/>
                    <GamesPage id="games"/>
                    <AdminPage id="admins"/>
                    <GroupsPage id="groups"/>
                </Epic>}
                {!loading && !error && !isAdmin && <RatingPage/>}
                {snack}
            </AppRoot>
        </GlobalActions.Provider>
    );
}

export default App;
