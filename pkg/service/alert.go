package service

/*
func (x *Service) insertAlert(ctx *types.Context, alert *model.Alert, url string, client db.Interface) (types.AlertID, error) {
	alert.Status = types.StatusNew
	alert.createdAt = time.Now().UTC()

	added, err := client.PutAlert(ctx, alert.toEnt())
	if err != nil {
		return "", err
	}

	if url != "" {
		endpoint := url + "/api/v1/alert/" + string(added.ID)
		alert.References = append(alert.References, &Reference{
			Title:   "alertchain API",
			URL:     endpoint,
			Comment: endpoint,
		})
	}

	if err := client.AddAttributes(ctx, added.ID, alert.Attributes.toEnt()); err != nil {
		return "", err
	}

	if err := client.AddReferences(ctx, added.ID, alert.References.toEnt()); err != nil {
		return "", err
	}

	return added.ID, nil
}

func (x *Service) CommitAlert(ctx *types.Context, id types.AlertID, req *model.ChangeRequest) error {
	if req.newStatus != nil {
		if err := client.UpdateAlertStatus(ctx, id, *req.newStatus); err != nil {
			return err
		}
	}
	if req.newSeverity != nil {
		if err := client.UpdateAlertSeverity(ctx, id, *req.newSeverity); err != nil {
			return err
		}
	}

	if len(req.newAttrs) > 0 {
		if err := client.AddAttributes(ctx, id, req.newAttrs); err != nil {
			return err
		}
	}

	for _, newAnn := range req.newAnnotations {
		if err := client.AddAnnotation(ctx, newAnn.attr, []*ent.Annotation{newAnn.ann}); err != nil {
			return err
		}
	}

	if len(req.newReferences) > 0 {
		if err := client.AddReferences(ctx, id, req.newReferences); err != nil {
			return err
		}
	}

	return nil
}
*/
