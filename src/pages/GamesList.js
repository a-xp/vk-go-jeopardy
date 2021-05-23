import {Button, CellButton, Div, Group, Panel, PanelHeader, PullToRefresh, SimpleCell} from "@vkontakte/vkui";
import {useCallback, useContext} from "react";
import {GlobalActions} from "../App";


function GamesList({onOpen, onRefresh, onNew, list, loading}) {

    const globalActions = useContext(GlobalActions);

    const onLinkCopy = useCallback((e) => {
        const id = e.currentTarget.getAttribute('data-id');
        globalActions.copyGameUrl(id);
        e.stopPropagation();
    }, [globalActions]);

    return (
        <Panel id="games">
            <PanelHeader>Викторины</PanelHeader>
            <Group>
                <PullToRefresh onRefresh={onRefresh} isFetching={loading}>
                    {list && list.map(game =>
                        <SimpleCell key={game.id}
                                    onClick={() => onOpen(game.id, 0)}
                                    description={game.id}
                                    after={<CellButton data-id={game.id} onClick={onLinkCopy}>Ссылка</CellButton>}
                                    indicator={game.active ? 'Активна' : (game.new ? 'Новая' : 'Архив')}>
                            {game.name}
                        </SimpleCell>
                    )}
                </PullToRefresh>
                <Div>
                    <Button stretched onClick={onNew}>Создать новую викторину</Button>
                </Div>
            </Group>
        </Panel>
    );

}

export default GamesList;
