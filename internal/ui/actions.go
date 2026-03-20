package ui

func (r *Renderer) goToMain() {
	r.currentView = mainView
	r.selectedTaskID = ""
	r.Render()
}

func (r *Renderer) selectTask(id string) {
	r.selectedTaskID = id
	r.currentView = taskActionView
	r.Render()
}

func (r *Renderer) Render() {
	r.clearMenu()

	switch r.currentView {
	case taskActionView:
		r.renderTaskActionMenu()
	default:
		r.renderMainMenu()
	}
}
