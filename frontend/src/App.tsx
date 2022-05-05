import React from 'react';
import './App.css';
import {QueryClient, QueryClientProvider, useMutation, useQuery, useQueryClient} from 'react-query'
import {ReactQueryDevtools} from 'react-query/devtools'
import {Button, Container, createStyles, Table, TextInput} from '@mantine/core';

import Client, {notes} from "./api/api";
import {useForm} from "@mantine/hooks";
import moment from "moment/moment";


const queryClient = new QueryClient()

var client = new Client(process.env.ENVIRONMENT || "local")

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

// {id: variables.id}], data


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
    const until_review = moment(note.next_review).fromNow()


    return (<tr key={note.id}>
        <td>{note.front}</td>
        <td>{note.back}</td>

        <td>{until_review}</td>
        <td>{moment(note.next_review).diff(moment(),'days')}</td>
        <td>
            {buttons.map(btn => {
                return <Button mr="xs" color={btn.color} compact={true}
                               onClick={() => handleClick(note.id, btn.answer)}>{btn.answer}</Button>
            })}

        </td>
    </tr>);
}

function Example() {

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

    if (error) return <div>'An error has occurred: ' + error</div>

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
            <td></td>
        </tr>
    )

    return (
        <form onSubmit={form.onSubmit((values) => newTodo(values))}>
            <Table>
                <thead>
                <tr>
                    <th>Front</th>
                    <th>Back</th>
                    <th>Next Review</th>
                    <th>Button</th>
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
