import * as React from "react";
import { useParams } from "react-router-dom";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";

import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";

import DnsIcon from "@material-ui/icons/Dns";
import SettingsEthernetIcon from "@material-ui/icons/SettingsEthernet";
import InsertDriveFileIcon from "@material-ui/icons/InsertDriveFile";
import PowerIcon from "@material-ui/icons/Power";
import LanguageIcon from "@material-ui/icons/Language";
import Rotate90DegreesCcwIcon from "@material-ui/icons/Rotate90DegreesCcw";
import ContactMailIcon from "@material-ui/icons/ContactMail";
import PersonIcon from "@material-ui/icons/Person";
import NoteIcon from "@material-ui/icons/Note";

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
            <Typography>
              Created at{" "}
              {strftime(
                "%Y-%m-%d %H:%M:%S",
                new Date(state.alert.created_at * 1000)
              )}
            </Typography>
          </Grid>

          <Grid>
            <List dense={true}>
              {state.alert.attributes ? (
                state.alert.attributes.map((attr) => {
                  const icon = attrIconMap[attr.type] || <NoteIcon />;
                  return (
                    <ListItem key={attr.id}>
                      <ListItemIcon>{icon}</ListItemIcon>
                      <ListItemText primary={attr.value} secondary={attr.key} />
                    </ListItem>
                  );
                })
              ) : (
                <div>No attributes</div>
              )}
            </List>
          </Grid>
        </Paper>
      </Grid>
    </div>
  );
}
