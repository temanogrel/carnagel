from cement.core.controller import CementBaseController


class PerformerController(CementBaseController):
    class Meta:
        label = 'performer'
        description = 'Performer related utilities'

        stacked_on = 'base'
        stacked_type = 'nested'




