import React from 'react';
import './App.css';
import {QueryClient, QueryClientProvider, useMutation, useQuery, useQueryClient} from 'react-query'
import {ReactQueryDevtools} from 'react-query/devtools'
import {Button, Container, createStyles, Table, TextInput} from '@mantine/core';

import Client, {notes} from "./api/api";
import {useForm} from "@mantine/hooks";
import moment from "moment/moment";


const queryClient = new QueryClient()
const environments: { [key: string]: string } = {
    "production": "production",
    "preview": "staging",
    "development": "staging",
    "local": "local",
}

let client = new Client(environments[process.env.REACT_APP_VERCEL_ENV || "local"])

const useMutateTodo = () => {
    const queryClient = useQueryClient()
    return useMutation(reviewNote, {
            onSuccess: (data, variables) => {
                queryClient.setQueryData<notes.ListResponse | undefined>("notes", (existingQuery) => {
                    if (existingQuery === undefined) {
                        return
                    }
                    return {
                        notes: existingQuery.notes.map(obj => {
                            return data.id === obj.id ? data : obj;
                        })
                    }
                });
            }
        }
    )
}

const useNewTodo = (p: () => void) => {
    const queryClient = useQueryClient()
    return useMutation(addNote, {
            onSuccess: (data, variables) => {
                queryClient.setQueryData<notes.ListResponse | undefined>("notes", (existingQuery) => {
                    if (existingQuery === undefined) {
                        return
                    }
                    return {
                        notes: [...existingQuery.notes, data]
                    }
                });
                p()
            }
        }
    )
}

interface reviewParams {
    id: string,
    answer: string
}

export const reviewNote = async ({id, answer}: reviewParams) => {
    return client.notes.ReviewNote(id, {answer: answer})
};

interface addParams {
    front: string,
    back: string
}

export const addNote = async ({front, back}: addParams) => {
    return client.notes.NewNote({front: front, back: back})
};

const buttons = [{
    color: "red",
    answer: "again",
    text: "-1"
}, {
    color: "gray",
    answer: "hard",
}, {
    color: "blue",
    answer: "good",
}, {
    color: "green",
    answer: "easy",
}]

const useStyles = createStyles((theme, _params, getRef) => ({
    reviewButton: {
        padding: theme.spacing.md,
        margin: theme.spacing.md,
    },
}));

function ReviewRow(note: notes.Note, handleClick: (id: string, answer: string) => void) {
    const until_review = moment(note.next_review).from(moment.utc())
    const is_due = moment(note.next_review).diff(moment(), 'seconds') < 0

    return (<tr key={note.id}>
        <td>{note.front}</td>
        <td>{note.back}</td>

        <td>
            {
                is_due ? buttons.map(btn => {
                    return <Button mr="xs" color={btn.color} compact={true} disabled={!is_due}
                                   onClick={() => handleClick(note.id, btn.answer)}>{btn.answer}</Button>
                }): "due " + until_review
            }

        </td>
    </tr>);
}

function Example() {
    console.log("starting in environment", process.env.ENVIRONMENT || "local")
    console.log(process.env)

    const form = useForm({
        initialValues: {
            front: '',
            back: '',
        },
    });

    const {isLoading, error, data} = useQuery<notes.ListResponse, Error>('notes', () => client.notes.ListNotes())
    const {mutate} = useMutateTodo();
    const {mutate: newTodo} = useNewTodo(() => form.reset());


    if (isLoading) return <div>'Loading...'</div>

    if (error) return <div>'An error has occurred: ' + <>{error}</></div>

    console.log(data)

    if (!data) {
        return null
    }


    const rows = data.notes.map((note) => {
            const handleClick = (id: string, answer: string) => {
                mutate({id: id, answer: answer});
            };
            return ReviewRow(note, handleClick);
        }
    );

    rows.push(
        <tr key="add">
            <td>
                <TextInput
                    placeholder="Front"
                    required
                    {...form.getInputProps('front')}
                />
            </td>
            <td>
                <TextInput
                    placeholder="Back"
                    required
                    {...form.getInputProps('back')}
                />
            </td>
            <td>
                <Button type="submit" color="dark" compact={true}>Add</Button>
            </td>
        </tr>
    )

    return (
        <form onSubmit={form.onSubmit((values) => newTodo(values))}>
            <Table>
                <thead>
                <tr>
                    <th>Front</th>
                    <th>Back</th>
                    <th>Review</th>
                </tr>
                </thead>
                <tbody>{rows}</tbody>
            </Table>
        </form>
    );
}

function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <Container size="md">
                <Example/>
            </Container>

            <ReactQueryDevtools initialIsOpen={false}/>
        </QueryClientProvider>
    );
}


export default App;
