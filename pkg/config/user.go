package config

type UserConfig struct {
	CustomResourcePanels []string
	History              *History `yaml:"history"`
}

func (c *UserConfig) AddCustomResourcePanels(resources ...string) {
	for _, resource := range resources {
		for _, each := range c.CustomResourcePanels {
			if each == resource {
				return
			}
		}
		c.CustomResourcePanels = append(c.CustomResourcePanels, resource)
	}
}

func (c *UserConfig) DeleteCustomResourcePanels(resources ...string) {
	for _, resource := range resources {
		for index, each := range c.CustomResourcePanels {
			if each == resource {
				c.CustomResourcePanels = append(c.CustomResourcePanels[:index], c.CustomResourcePanels[index+1:]...)
			}
		}
	}
}

type History struct {
	ImageHistory   []string `yaml:"image_history"`
	CommandHistory []string `yaml:"command_history"`
	PodNameHistory []string `yaml:"pod_name_history"`
}

func (h *History) AddStringHistory(history []string, newOne string) []string {
	history = append([]string{newOne}, history...)
	for index, each := range history[1:] {
		if each == newOne {
			history = append(history[:index], history[index+1:]...)
		}
	}
	return history
}

func (h *History) AddCommandHistory(newOne string) {
	h.CommandHistory = h.AddStringHistory(h.CommandHistory, newOne)
}

func (h *History) AddImageHistory(newOne string) {
	h.ImageHistory = h.AddStringHistory(h.ImageHistory, newOne)
}

func (h *History) AddPodNameHistory(newOne string) {
	h.PodNameHistory = h.AddStringHistory(h.PodNameHistory, newOne)
}
