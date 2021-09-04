import * as React from "react";
import { useParams } from "react-router-dom";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import Link from "@material-ui/core/Link";
import LinkIcon from "@material-ui/icons/Link";

import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import Button from "@material-ui/core/Button";

import DnsIcon from "@material-ui/icons/Dns";
import SettingsEthernetIcon from "@material-ui/icons/SettingsEthernet";
import InsertDriveFileIcon from "@material-ui/icons/InsertDriveFile";
import PowerIcon from "@material-ui/icons/Power";
import LanguageIcon from "@material-ui/icons/Language";
import Rotate90DegreesCcwIcon from "@material-ui/icons/Rotate90DegreesCcw";
import ContactMailIcon from "@material-ui/icons/ContactMail";
import PersonIcon from "@material-ui/icons/Person";
import NoteIcon from "@material-ui/icons/Note";

import ErrorOutlineOutlinedIcon from "@material-ui/icons/ErrorOutlineOutlined";
import CheckBoxIcon from "@material-ui/icons/CheckBox";
import CircularProgress from "@material-ui/core/CircularProgress";

import strftime from "strftime";

import Alert from "@material-ui/lab/Alert";

import * as model from "./model";

type alertProps = {
  id?: string;
};

type alertState = {
  isLoaded: boolean;
  alert?: model.alert;
  error?: string;
};

const attrIconMap = {
  ipaddr: <SettingsEthernetIcon />,
  domain: <DnsIcon />,
  port: <PowerIcon />,
  userid: <PersonIcon />,
  email: <ContactMailIcon />,
  sha256: <Rotate90DegreesCcwIcon />,
  filepath: <InsertDriveFileIcon />,
  url: <LanguageIcon />,
};

const taskStatusIconMap = {
  succeeded: <CheckBoxIcon />,
  failure: <ErrorOutlineOutlinedIcon />,
  running: <CircularProgress />,
};

export function AlertView(props: alertProps) {
  const { alertID } = useParams();
  const id = props.id ? props.id : alertID;

  const [state, setState] = React.useState<alertState>({
    isLoaded: false,
  });

  const getAlert = () => {
    fetch(`/api/v1/alert/` + id)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log({ result });
          if (result.data) {
            setState({
              isLoaded: true,
              alert: result.data as model.alert,
            });
          } else {
            setState({
              isLoaded: true,
            });
          }
        },
        (error) => {
          console.log({ error });
          setState({
            isLoaded: true,
            error: error.message,
          });
        }
      );
  };

  React.useEffect(getAlert, [id]);

  if (!state.isLoaded) {
    return <div>Loading...</div>;
  }
  if (state.error) {
    return <Alert severity="error">{state.error}</Alert>;
  }

  return (
    <div>
      <Grid>
        <Typography variant="h1">{state.alert.title}</Typography>
      </Grid>
      <Grid>
        <Paper>
          <Grid>
            <Typography variant="h5">Summary</Typography>
          </Grid>
          <Grid>
            <List>
              <ListItem>Description: {state.alert.description}</ListItem>
              <ListItem>Status: {state.alert.status}</ListItem>
              <ListItem>Severity: {state.alert.severity}</ListItem>
              <ListItem>Detected by: {state.alert.detector}</ListItem>
              <ListItem>
                Created at:
                {strftime(
                  " %Y-%m-%d %H:%M:%S",
                  new Date(state.alert.created_at * 1000)
                )}
              </ListItem>
              <ListItem>
                Closed at:
                {state.alert.closed_at
                  ? strftime(
                      " %Y-%m-%d %H:%M:%S",
                      new Date(state.alert.closed_at * 1000)
                    )
                  : "N/A"}
              </ListItem>
            </List>
          </Grid>
        </Paper>
      </Grid>

      <Grid>
        <Paper>
          <Grid>
            <Typography variant="h5">References</Typography>
          </Grid>
          <Grid>
            <List dense={true}>
              {state.alert.references.map((ref) => {
                return (
                  <ListItem key={ref.id}>
                    <ListItemIcon>
                      <LinkIcon />
                    </ListItemIcon>
                    <Link href={ref.url}>
                      <ListItemText
                        primary={ref.title}
                        secondary={ref.comment}
                      />
                    </Link>
                  </ListItem>
                );
              })}
            </List>
          </Grid>
        </Paper>
      </Grid>

      <Grid>
        <Attributes attrs={state.alert.attributes} />
      </Grid>

      <Grid>
        <Paper>
          <Grid>
            <Typography variant="h5">Tasks</Typography>
          </Grid>

          {state.alert.task_logs.map((task_log) => {
            return (
              <Grid key={task_log.id} container>
                <Grid item>{taskStatusIconMap[task_log.status]}</Grid>
                <Grid item>
                  <Typography variant="h6">{task_log.task_name}</Typography>
                </Grid>
              </Grid>
            );
          })}
        </Paper>
      </Grid>
    </div>
  );
}

type AttributesProps = {
  attrs: model.attribute[];
};

function Attributes(props: AttributesProps) {
  const execAction = (actionID: string, attrID: number) => {
    const data = { action_id: actionID, attr_id: attrID };
    fetch(`/api/v1/action/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    })
      .then((res) => res.json())
      .then(
        (result) => {
          console.log({ result });
        },
        (error) => {
          console.log({ error });
        }
      );
  };

  return (
    <Paper>
      <Grid>
        <Typography variant="h5">Attributes</Typography>
      </Grid>

      <Grid>
        <List dense={true}>
          {props.attrs ? (
            props.attrs.map((attr) => {
              const icon = attrIconMap[attr.type] || <NoteIcon />;
              return (
                <ListItem key={attr.id} style={{ width: 500 }}>
                  <ListItemIcon>{icon}</ListItemIcon>
                  <ListItemText primary={attr.value} secondary={attr.key} />
                  <div>
                    {(attr.actions || []).map((action, idx) => {
                      return (
                        <Button
                          key={idx}
                          variant="contained"
                          onClick={() => {
                            execAction(action.id, attr.id);
                          }}
                          color="primary">
                          {action.name}
                        </Button>
                      );
                    })}
                  </div>
                </ListItem>
              );
            })
          ) : (
            <div>No attributes</div>
          )}
        </List>
      </Grid>
    </Paper>
  );
}
