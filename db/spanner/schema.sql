CREATE TABLE Alerts (
    ID STRING(MAX) NOT NULL,
    Title STRING(MAX) NOT NULL,
    Description STRING(MAX),
    Detector STRING(MAX) NOT NULL,
    Severity STRING(MAX) NOT NULL,

    DetectedAt TIMESTAMP,

    CreatedAt TIMESTAMP NOT NULL,
    ClosedAt TIMESTAMP,
) PRIMARY KEY(ID);
CREATE INDEX AlertByCreatedAt ON Alerts (CreatedAt);
CREATE INDEX AlertByClosedAt ON Alerts (ClosedAt);

CREATE TABLE Attributes (
    ID STRING(MAX) NOT NULL,
    AlertID STRING(MAX) NOT NULL,

    Key STRING(MAX) NOT NULL,
    Value STRING(MAX) NOT NULL,
    Type STRING(MAX) NOT NULL,
    Contexts ARRAY<STRING(MAX)>,

    CONSTRAINT FK_Attributes_AlertID FOREIGN KEY (AlertID) REFERENCES Alerts (ID),
) PRIMARY KEY(ID);
CREATE INDEX AttributesByAlertID ON Attributes (AlertID);

CREATE TABLE Annotations (
    ID STRING(MAX) NOT NULL,
    AlertID STRING(MAX) NOT NULL,
    AttributeID STRING(MAX) NOT NULL,

    Timestamp TIMESTAMP,
    Source STRING(MAX) NOT NULL,
    Name STRING(MAX) NOT NULL,
    Value STRING(MAX) NOT NULL,
    Tags ARRAY<STRING(MAX)>,
    URI STRING(MAX),

    CONSTRAINT FK_Annotation_AlertID FOREIGN KEY (AlertID) REFERENCES Alerts (ID),
    CONSTRAINT FK_Annotation_AttributeID FOREIGN KEY (AttributeID) REFERENCES Attributes (ID),
) PRIMARY KEY(ID);
CREATE INDEX AnnotationsByAlertID ON Annotations (AlertID);
CREATE INDEX AnnotationsByAttributeID ON Annotations (AttributeID);

CREATE TABLE References (
    ID STRING(MAX) NOT NULL,
    AlertID STRING(MAX) NOT NULL,

    Source STRING(MAX) NOT NULL,
    Title STRING(MAX) NOT NULL,
    URI STRING(MAX) NOT NULL,
    Comment STRING(MAX),

    CONSTRAINT FK_Reference_AlertID FOREIGN KEY (AlertID) REFERENCES Alerts (ID),
) PRIMARY KEY(ID);
CREATE INDEX ReferencesByAlertID ON References (AlertID);
