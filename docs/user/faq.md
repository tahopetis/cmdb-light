# CMDB Lite Frequently Asked Questions (FAQ)

This FAQ addresses common questions and issues that users may encounter while using CMDB Lite. If you don't find an answer to your question here, please check the [User Guide](README.md) or contact your system administrator.

## General Questions

### What is CMDB Lite?

CMDB Lite is a lightweight Configuration Management Database (CMDB) application that helps IT teams store, track, and visualize Configuration Items (CIs) and their relationships. It provides a simple, intuitive interface for managing infrastructure components and their dependencies.

### Who is CMDB Lite designed for?

CMDB Lite is designed for small to mid-sized organizations or teams that need visibility and auditability over their infrastructure without the complexity of a full ITSM suite. It's particularly useful for IT Operations, DevOps/SRE teams, and Security/Audit teams.

### What are the system requirements for using CMDB Lite?

CMDB Lite is a web-based application, so you only need a modern web browser (Chrome, Firefox, Safari, or Edge) to use it. The system requirements for hosting CMDB Lite are detailed in the [Deployment Guide](../operator/deployment.md).

### Is there a mobile app for CMDB Lite?

Currently, CMDB Lite doesn't have a dedicated mobile app, but the web interface is responsive and works well on mobile devices. We're considering developing a mobile app in the future.

## Account and Access

### How do I get access to CMDB Lite?

Access to CMDB Lite is typically managed by your system administrator. If you need access, please contact your IT department or the person responsible for managing CMDB Lite in your organization.

### I forgot my password. How can I reset it?

Currently, password reset must be handled by your system administrator. Contact them to reset your password. We're working on implementing a self-service password reset feature in a future release.

### What are the different user roles in CMDB Lite?

CMDB Lite has three predefined user roles:

- **Admin**: Full access to all features and settings, including user management
- **Editor**: Can create, update, and delete CIs and relationships, but cannot manage users or system settings
- **Viewer**: Read-only access to CIs and relationships, cannot make changes

### How do I change my role in CMDB Lite?

User roles can only be changed by administrators. If you believe you need a different role, please contact your system administrator.

## Configuration Items (CIs)

### What types of CIs can I create in CMDB Lite?

CMDB Lite supports several predefined CI types, including Servers, Applications, Databases, Network Devices, and Storage Systems. You can also create custom CI types to meet your specific needs.

### How do I create a new CI type?

Currently, creating new CI types requires administrator privileges. Contact your system administrator if you need a new CI type. We're working on making this more flexible in future releases.

### What attributes can I add to a CI?

You can add custom attributes to any CI as key-value pairs. These can include text, numbers, dates, and even JSON objects for complex data. There are no restrictions on the types of attributes you can add.

### How do I delete a CI?

To delete a CI:

1. Navigate to the CI detail view by clicking on the CI name in the CI list.
2. Click the "Delete" button.
3. Confirm the deletion in the popup dialog.

Note: Deleting a CI will also remove all relationships associated with it, so be sure this is what you want to do.

### Can I recover a deleted CI?

Currently, there is no built-in way to recover deleted CIs. However, your system administrator may be able to restore from a backup. We're planning to implement a "soft delete" feature in the future that would allow recovery of deleted items.

## Relationships

### What types of relationships can I create between CIs?

CMDB Lite supports several predefined relationship types, including "Depends On", "Hosts", "Connects To", and "Part Of". You can also create custom relationship types to model your specific infrastructure dependencies.

### How do I create a relationship between two CIs?

To create a relationship:

1. Navigate to the detail view of the source CI.
2. Click the "Add Relationship" button.
3. Select the target CI from the dropdown list.
4. Choose the relationship type.
5. Click "Save" to create the relationship.

### Can I create circular relationships (A depends on B, B depends on A)?

Yes, CMDB Lite allows you to create circular relationships. However, you should be careful when doing so, as it can make impact analysis more complex.

### How do I visualize relationships between CIs?

CMDB Lite provides a graph visualization feature that shows CIs and their relationships as an interactive network diagram. You can access this by clicking on the "Graph" link in the main navigation.

## Graph Visualization

### Why is the graph visualization slow or unresponsive?

The graph visualization may become slow if you have a large number of CIs and relationships. To improve performance:

1. Use the filters to show only specific types of CIs.
2. Use the search function to find and focus on specific CIs.
3. Try collapsing nodes that you're not currently interested in.

We're continuously working on improving the performance of the graph visualization.

### Can I export the graph as an image?

Currently, there is no built-in way to export the graph as an image. This is a feature we're planning to add in a future release. In the meantime, you can use your operating system's screenshot functionality to capture the graph.

### Why don't all my CIs appear in the graph?

The graph visualization may not show all CIs if:

1. You have filters applied that exclude certain CI types.
2. The CIs have no relationships to other CIs and are filtered out for clarity.
3. There are too many CIs to display at once, and some are hidden for performance reasons.

You can adjust the filters and use the search function to find specific CIs in the graph.

## Search and Filtering

### How do I search for CIs?

You can search for CIs using the search bar at the top of the interface. Simply type what you're looking for, and the results will be filtered in real-time. The search looks at CI names, attributes, and tags.

### Can I save my searches or filters?

Currently, there is no built-in way to save searches or filters. This is a feature we're planning to add in a future release. In the meantime, you can bookmark specific filter URLs in your browser.

### How do I search within specific attributes?

The search function looks across all attributes by default. To search within specific attributes, you can use the advanced search syntax:

- `attribute:value` - Search for CIs where the attribute contains the specified value
- `tag:tagname` - Search for CIs with a specific tag
- `type:typename` - Search for CIs of a specific type

## Audit Logs

### How long are audit logs retained?

Audit log retention is configured by your system administrator. Typically, logs are retained for a period of 1-2 years, but this can vary based on your organization's policies.

### Can I export audit logs?

Yes, you can export audit logs for external analysis or compliance reporting. Use the export function in the audit log view to download logs in CSV or JSON format.

### Why don't I see all actions in the audit logs?

Audit logs should capture all significant actions in the system. If you believe something is missing, it could be due to:

1. Filtering that is hiding certain types of actions
2. A time range that doesn't include the period when the action occurred
3. A potential issue with the logging system

Contact your system administrator if you suspect an issue with audit logging.

## Integration and API

### Does CMDB Lite have an API?

Yes, CMDB Lite provides a RESTful API that allows programmatic access to all features. The API uses JWT tokens for authentication and follows RESTful conventions.

### How do I get API access?

API access requires authentication with a JWT token. You can generate API tokens in your user settings, or your system administrator can provide you with an API key. Contact your administrator for more information.

### Are there any rate limits for API calls?

Yes, the API has configurable rate limits to prevent abuse. The default limits are set to allow reasonable usage while preventing potential denial-of-service attacks. If you need higher limits, contact your system administrator.

### Can I integrate CMDB Lite with other tools?

Yes, CMDB Lite is designed to be integrable with other tools. You can use the API to integrate with monitoring systems, configuration management tools, and other IT operations tools. We also support webhooks for event notifications (future feature).

## Troubleshooting

### I'm having trouble logging in. What should I do?

If you're having trouble logging in:

1. Verify that you're using the correct username and password.
2. Check if your account is locked or disabled.
3. Try clearing your browser cache and cookies.
4. Try using a different browser or incognito/private mode.
5. Contact your system administrator if the issue persists.

### The interface is running slowly. How can I improve performance?

If the interface is slow:

1. Check your internet connection.
2. Try closing unnecessary browser tabs.
3. Clear your browser cache.
4. Use filters to reduce the amount of data being displayed.
5. Contact your system administrator if the issue persists, as it may be a server-side problem.

### I'm getting an error message. What does it mean?

CMDB Lite provides descriptive error messages to help you understand what went wrong. Common errors include:

- **Authentication Error**: Your login credentials are incorrect or your session has expired.
- **Permission Denied**: You don't have permission to perform the requested action.
- **Validation Error**: The data you entered is not valid (e.g., missing required fields).
- **Server Error**: There was a problem on the server. Try again later or contact your administrator.

If you don't understand an error message or need help resolving it, please contact your system administrator.

## Data Management

### How do I back up my CMDB data?

Data backup is typically handled by your system administrator. If you need to create your own backups, contact your administrator for guidance on the backup procedures.

### Can I import data from another CMDB or spreadsheet?

Currently, CMDB Lite doesn't have built-in import functionality, but we're planning to add CSV and JSON import features in a future release. In the meantime, you can use the API to import data from other sources.

### How do I export my CMDB data?

You can export CI data using the export function in the CI list view. This allows you to download data in CSV or JSON format for external analysis or reporting.

## Feature Requests and Feedback

### How do I request a new feature?

We welcome feature requests! To request a new feature:

1. Check if the feature has already been requested in our [GitHub Issues](https://github.com/yourorg/cmdb-lite/issues).
2. If not, create a new issue with the "feature request" label.
3. Provide a clear description of the feature and its benefits.
4. Explain how it would fit into your workflow.

### How do I report a bug?

To report a bug:

1. Check if the bug has already been reported in our [GitHub Issues](https://github.com/yourorg/cmdb-lite/issues).
2. If not, create a new issue with the "bug" label.
3. Provide a detailed description of the problem, including steps to reproduce it.
4. Include any error messages or screenshots that might help diagnose the issue.

### How can I provide feedback on CMDB Lite?

We value your feedback! You can provide feedback by:

1. Creating an issue in our [GitHub Issues](https://github.com/yourorg/cmdb-lite/issues).
2. Contacting your system administrator, who can relay feedback to the development team.
3. Participating in community discussions (future feature).

## Still Have Questions?

If you didn't find an answer to your question in this FAQ, please:

1. Check the [User Guide](README.md) for more detailed instructions.
2. Contact your system administrator for account or access issues.
3. Search our [GitHub Issues](https://github.com/yourorg/cmdb-lite/issues) for similar questions.
4. Create a new issue with your question if it hasn't been addressed before.