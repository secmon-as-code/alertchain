import * as React from "react";
import * as ReactDOM from "react-dom";
import {
  createTheme,
  createStyles,
  ThemeProvider,
  makeStyles,
  Theme,
} from "@material-ui/core/styles";
import CssBaseline from "@material-ui/core/CssBaseline";
import Typography from "@material-ui/core/Typography";
import Link from "@material-ui/core/Link";
import Box from "@material-ui/core/Box";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Button from "@material-ui/core/Button";
import IconButton from "@material-ui/core/IconButton";
import MenuIcon from "@material-ui/icons/Menu";

import { AlertList } from "./components/alertList";
import { AlertView } from "./components/alert";

import {
  BrowserRouter,
  Route,
  Switch,
  Redirect,
  Link as RouterLink,
} from "react-router-dom";

function App() {
  const classes = useStyle();

  return (
    <ThemeProvider theme={theme}>
      <div className={classes.root}>
        <Box sx={{ flexGrow: 1 }}>
          <AppBar position="static">
            <Toolbar>
              <Typography variant="h5">AlertChain</Typography>
            </Toolbar>
          </AppBar>
          <main className={classes.main}>
            <BrowserRouter>
              <CssBaseline />
              <div className={classes.app}>
                <main className={classes.main}>
                  <Switch>
                    <Route path="/alert/:alertID">
                      <AlertView />
                    </Route>
                    <Route path="/alert">
                      <AlertList />
                    </Route>
                    <Route path="/" exact>
                      <Redirect to="/alert" />
                    </Route>
                  </Switch>
                </main>
                <footer className={classes.footer}>
                  <Copyright />
                </footer>
              </div>
            </BrowserRouter>
          </main>
        </Box>
      </div>
    </ThemeProvider>
  );
}

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {"Copyright © "}
      <Link color="inherit" href="https://github.com/m-mizutani">
        m-mizutani
      </Link>{" "}
      {new Date().getFullYear()}
      {"."}
    </Typography>
  );
}

let theme = createTheme({
  palette: {
    primary: {
      light: "#757ce8",
      main: "#3f50b5",
      dark: "#002884",
      contrastText: "#fff",
    },
  },

  typography: {
    h1: {
      fontWeight: "bold",
      fontSize: 48,
      letterSpacing: 0.5,
      //   fontFamily: ["Kanit"].join(","),
    },
    h5: {
      fontWeight: "bold",
      fontSize: 20,
      letterSpacing: 0.1,
    },
  },
});

theme = {
  ...theme,
  overrides: {
    MuiDrawer: {
      paper: {
        backgroundColor: "#18202c",
      },
    },
    MuiButton: {
      label: {
        textTransform: "none",
      },
      contained: {
        boxShadow: "none",
        "&:active": {
          boxShadow: "none",
        },
      },
    },
    MuiTabs: {
      root: {
        marginLeft: theme.spacing(1),
      },
      indicator: {
        height: 3,
        borderTopLeftRadius: 3,
        borderTopRightRadius: 3,
        backgroundColor: theme.palette.common.white,
      },
    },
    MuiTab: {
      root: {
        textTransform: "none",
        margin: "0 16px",
        minWidth: 0,
        padding: 0,
        [theme.breakpoints.up("md")]: {
          padding: 0,
          minWidth: 0,
        },
      },
    },
    MuiIconButton: {
      root: {
        padding: theme.spacing(1),
      },
    },
    MuiTooltip: {
      tooltip: {
        borderRadius: 4,
      },
    },
    MuiDivider: {
      root: {
        backgroundColor: "#404854",
      },
    },
    MuiListItemText: {
      primary: {
        fontWeight: theme.typography.fontWeightMedium,
      },
    },
    MuiListItemIcon: {
      root: {
        color: "inherit",
        marginRight: 0,
        "& svg": {
          fontSize: 20,
        },
      },
    },
    MuiAvatar: {
      root: {
        width: 32,
        height: 32,
      },
    },
  },
};

const useStyle = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      display: "flex",
      minHeight: "100vh",
    },
    paper: {
      margin: "auto",
      overflow: "hidden",
    },
    contentWrapper: {
      margin: "40px 30px",
    },
    app: {
      flex: 1,
      display: "flex",
      flexDirection: "column",
    },
    main: {
      flex: 1,
      padding: theme.spacing(6, 4),
      background: "#fff",
    },
    footer: {
      padding: theme.spacing(2),
      background: "#fff",
    },
  })
);

ReactDOM.render(<App />, document.querySelector("#app"));
