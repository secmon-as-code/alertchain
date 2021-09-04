import * as React from "react";
import Grid from "@material-ui/core/Grid";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import { Link as RouterLink } from "react-router-dom";

import strftime from "strftime";
import * as model from "./model";

const useStyle = makeStyles((theme: Theme) =>
  createStyles({
    alertItem: {
      fontSize: 18,
      margin: "10px 20px",
    },
  })
);

type alertListState = {
  isLoaded: boolean;
  alerts?: model.alert[];
  error?: string;
};

export function AlertList() {
  return (
    <div>
      <Grid>
        <Typography variant="h1">Alert</Typography>
      </Grid>
      <Grid>
        <ListView />
      </Grid>
    </div>
  );
}

function ListView() {
  const classes = useStyle();

  const [state, setState] = React.useState<alertListState>({
    isLoaded: false,
  });

  const getAlertList = () => {
    fetch(`/api/v1/alert`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log({ result });
          if (result.data) {
            setState({
              isLoaded: true,
              alerts: result.data as model.alert[],
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

  React.useEffect(getAlertList, []);

  if (!state.isLoaded) {
    return <div>Loading</div>;
  }
  if (!state.alerts) {
    return <div>No alerts</div>;
  }

  return (
    <div style={{ margin: 20 }}>
      {state.alerts.map((alert) => {
        return (
          <Grid
            container
            key={alert.id}
            style={{
              marginTop: -1,
              border: "1px solid #bbb",
            }}>
            <Grid item className={classes.alertItem} style={{ minWidth: 600 }}>
              <Grid style={{ margin: 5 }}>
                <RouterLink to={"/alert/" + alert.id}>
                  <Typography variant="h5">{alert.title}</Typography>
                </RouterLink>
              </Grid>
              <Grid style={{ fontSize: 14, color: "#444" }}>
                detected by {alert.detector} at{" "}
                {strftime(
                  "%Y-%m-%d %H:%M:%S",
                  new Date(alert.created_at * 1000)
                )}
              </Grid>
            </Grid>
            <Grid item className={classes.alertItem} style={{ minWidth: 100 }}>
              {alert.status}
            </Grid>
            <Grid item className={classes.alertItem} style={{ minWidth: 100 }}>
              {alert.severity || "N/A"}
            </Grid>
          </Grid>
        );
      })}
    </div>
  );
}
