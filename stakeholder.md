# ideas
- export cv as markdown file
- export as full bundle (pdf, docx, json, markdown)
- we currently support standard cv layouts. We might want to support an entire new class of CV called the consultant profile. We have receive a template from one of our stakeholders for cvs that are used to present skilled professional consultants from various fields to potential clients. In these the structure is different since the main point of focus are skills, past projects and certifications
- a feature to upload an existing cv (docx or pdf) and have the informations we need to build our cv extracted and set up as a new cv. this would simplify onboarding of customers that do not have a markdown or json base for their current cvs
- a feature to actually import back the informations from a cv we have previously extracted as json or markdown. The feature should be solid enough to be able to also handle modifications made to the structure or contents so we can get a result regardless of the user having edited the file
  

# Bugs
- it seems we still have bugs saving a cv after we made changes to it. We suspect the switch of the templates makes the saving get broken in the backend
- the frontend is still not responsive enough, we should do way more testing and adjusting in terms of viewports and also different clients. We want to have mobile users have also a good time editing cvs.
-  Much clipping issues in regards to half-screen windows, looks really off

# Findings
- currently each clicked section will have a button displayed that says "click to edit". The button is still there even when you are already in the opened sidebar or modal to edit. It should either be removed when you are in the edit pane or not displayed at all
- when clicking any section you often are not directly able to edit them. We should build our components in a way, that when you click them in a subcomponents (e.g. textfield in section) you should be able to edit it directly without an extra click (maybe keyboard focus directly in the edit field needed)
- the frontend looks really bland and boring. We need to do a full design overhaul and generate our visual direction for the app
- the general export of CVs in DOCX works but the styling is mismatched and does not look like the actual CV design. We need to update it to match. This task should be tackled after we have finalized the stylings of the CV templates
- the date pickers are very uggly and also not good for the user experience, we should use some that are better
- a free text field for technologies is not what one should use. If we say we have one entry per line we should build something that tokenizes the entries from the get go. 
- our documentation for the repository is very outdated and also super bland. We really need to update our readme in terms of content and also design. lets use the following readme as inspiration: https://github.com/openclaw/openclaw/blob/main/README.md